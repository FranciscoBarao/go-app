package database

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/FranciscoBarao/catalog/config"
	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/middleware/logging"
	"github.com/FranciscoBarao/catalog/model"
)

type Postgres struct {
	db *gorm.DB
}

func Connect(config *config.PostgresConfig) (*Postgres, error) {
	log := logging.FromCtx(context.Background())

	db, err := gorm.Open(postgres.Open(config.String()), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return nil, err
	}

	log.Debug().Msg("connected to database")

	if err = migrate(db, &model.Boardgame{}); err != nil {
		return nil, err
	}
	if err = migrate(db, &model.Tag{}); err != nil {
		return nil, err
	}
	if err = migrate(db, &model.Category{}); err != nil {
		return nil, err
	}
	if err = migrate(db, &model.Mechanism{}); err != nil {
		return nil, err
	}
	if err = migrate(db, &model.Rating{}); err != nil {
		return nil, err
	}

	log.Debug().Msg("database migration completed")

	return &Postgres{db}, nil
}

func migrate(db *gorm.DB, model interface{}) error {
	log := logging.FromCtx(context.Background())
	err := db.AutoMigrate(model)
	if err != nil {
		log.Error().Err(err).Interface("model", model).Msg("failed to migrate model")
	}
	return err
}

func isSliceOrArray(value interface{}) bool {
	return reflect.ValueOf(value).Elem().Kind() == reflect.Slice || reflect.ValueOf(value).Elem().Kind() == reflect.Array
}

func (instance *Postgres) Create(value interface{}, omits ...string) error {
	log := logging.FromCtx(context.Background())

	err := instance.db.Omit(omits...).Create(value).Error
	if err != nil {
		log.Error().Err(err).Interface("value", value).Msg("failed to create database entry")
		if errors.Is(err, gorm.ErrRegistered) {
			return middleware.NewError(http.StatusConflict, "Entry already registered")
		}
		return err
	}

	log.Debug().Interface("value", value).Msg("created database entry")
	return nil
}

func (instance *Postgres) Read(value interface{}, sort, search, identifier string) error {
	log := logging.FromCtx(context.Background())

	var err error
	if isSliceOrArray(value) {
		if search == "" {
			err = instance.db.Preload(clause.Associations).Order(sort).Find(value).Error // Find all with sort and NO filters
		} else {
			err = instance.db.Preload(clause.Associations).Order(sort).Find(value, search, identifier).Error // Find all with filters and sort
		}
	} else {
		err = instance.db.Preload(clause.Associations).First(value, search, identifier).Error // Find 1 Specific
	}

	if err != nil {
		log.Error().Err(err).Str("search", search).Str("identifier", identifier).Msg("failed to read database entry")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return middleware.NewError(http.StatusNotFound, "Record not found")
		}
		return err
	}

	log.Debug().Interface("value", value).Msg("fetched database entry")
	return nil
}

func (instance *Postgres) Update(value interface{}, omits ...string) error {
	log := logging.FromCtx(context.Background())

	err := instance.db.Omit(omits...).Save(value).Error
	if err != nil {
		log.Error().Err(err).Interface("value", value).Msg("failed to update database entry")
		return err
	}

	log.Debug().Interface("value", value).Msg("updated database entry")
	return nil
}

func (instance *Postgres) Delete(value interface{}) error {
	log := logging.FromCtx(context.Background())

	// Delete BG and all its associations (E.g Tags associations)
	err := instance.db.Select(clause.Associations).Delete(value).Error
	if err != nil {
		log.Error().Err(err).Interface("value", value).Msg("failed to delete database entry")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return middleware.NewError(http.StatusNotFound, "Record Not found")
		}
		return err
	}
	log.Debug().Interface("value", value).Msg("deleted database entry")
	return nil
}

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<        ASSOCIATIONS        >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// AppendAssociatons adds certain associations to a certain model (E.g Add Tags to a Boardgame)
func (instance *Postgres) AppendAssociatons(model interface{}, association string, values interface{}) error {
	log := logging.FromCtx(context.Background())

	err := instance.db.Model(model).Association(association).Append(values)
	if err != nil {
		log.Error().Err(err).Interface("model", model).Str("association", association).Interface("values", values).Msg("failed to append associations")
		return err
	}

	log.Debug().Interface("model", model).Str("association", association).Interface("values", values).Msg("associated")
	return nil
}

// ReadAssociatons gets associations of a type of a certain model (E.g Get Tags of a Boardgame)
func (instance *Postgres) ReadAssociatons(model interface{}, association string, store interface{}) error {
	log := logging.FromCtx(context.Background())

	err := instance.db.Model(model).Association(association).Find(store)
	if err != nil {
		log.Error().Err(err).Interface("model", model).Str("association", association).Msg("failed to read associations")
		return err
	}

	log.Debug().Interface("model", model).Str("association", association).Interface("store", store).Msg("fetched associations")
	return nil
}

// ReplaceAssociatons replaces the values of a certain association of a certain model (E.g Replace Tags of a Boardgame)
func (instance *Postgres) ReplaceAssociatons(model interface{}, association string, values interface{}) error {
	log := logging.FromCtx(context.Background())

	err := instance.db.Model(model).Association(association).Replace(values)
	if err != nil {
		log.Error().Err(err).Interface("model", model).Str("association", association).Interface("values", values).Msg("failed to replace associations")
		return err
	}

	log.Debug().Interface("model", model).Str("association", association).Interface("values", values).Msg("replaced association")
	return nil
}

// MDeleteAssociatons deletes all values of a certain association of a certain model (E.g Delete all Tags of a Boardgame)
func (instance *Postgres) DeleteAssociatons(model interface{}, association string) error {
	log := logging.FromCtx(context.Background())

	err := instance.db.Model(model).Association(association).Clear()
	if err != nil {
		log.Error().Err(err).Interface("model", model).Str("association", association).Msg("failed to delete associations")
		return err
	}

	log.Debug().Interface("model", model).Str("association", association).Msg("deleted associations")
	return nil
}
