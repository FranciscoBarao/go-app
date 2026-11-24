# Service Name

<!-- Start the introduction with "Display Name (`repo-name`)".
     Explain what the service does, what problem it solves, and who consumes it in 3-5 sentences. -->

## Design

<!-- Describe the high-level design and the service's place in the platform.
     Include a diagram when it improves understanding, and link to relevant design documents. -->

## Dependencies

<!-- Identify the services and infrastructure this service calls.
     Explain what each dependency provides and in which flow it is used.
     Include external infrastructure such as PostgresSQL where applicable. -->

## Consumers

<!-- Identify each service or external actor that calls this service and why.
     Distinguish game clients, game servers, operators, controllers, and internal services.
     Summarize whether each interface is externally routed, cluster-internal, or pod-local. -->

## Data Model

<!-- Remove this section if not applicable. -->
<!-- Describe the data this service uniquely owns.
     Include key model definitions and link to the API contract or implementation. -->

### Model Validations

<!-- Remove this section if not applicable. -->
<!-- Describe non-obvious domain, configuration, or platform invariants enforced before data is persisted,
     deployed, or executed. Explain why each constraint exists, how failures are reported, and link to its
     implementation. Do not repeat routine required-field checks already clear from the API contract. -->

## Data Store

<!-- Remove this section if not applicable. -->
<!-- Identify the data store and distinguish persistent, cached, and ephemeral state.
     State which component owns each write path. -->

## Data / Process Flows

<!-- Remove this section if not applicable. -->
<!-- Illustrate the critical data or process flows -->

## Interfaces

<!-- For every interface, record its audience, exposure path, authentication, authorization, and enforcement point.
     Distinguish TLS or network containment from application controls. State "no application-level control verified"
     or "evidence not found at the inspected revision" when applicable instead of inferring a control. -->

### HTTP API

<!-- Link to definitions in local API files and describe the main endpoints.
     Identify public and private servers separately, including any method-specific access policy.
     Document deliberate unauthenticated endpoints -->

### Event Bus

<!-- Remove this section if not applicable. -->
<!-- Describe events published or consumed through a bus.
     Record what triggers each event, the payload ownership, and delivery or retry guarantees. -->

### Exports

<!-- Remove this section if not applicable. -->
<!-- Describe exported packages used by other repositories and link to their documentation. -->

# Contributing

## Testing
<!-- Document how to run unit and integration tests. Use `make check` when available 
Do not instruct contributors to edit generated mocks. Mock generation remains repository or CI owned. -->

```shell
make check
```

## Running Locally

<!-- Remove this section if not applicable. -->
<!-- Describe supported local startup, required dependencies, and safe configuration.
     Do not copy credentials or environment-specific secret values into the README. -->

## Building

<!-- Document the build command supported by this repository. 
     Do not assume every service exposes the same Makefile targets. -->


# Operations

## Deployment

<!-- Describe workload type, deployment target, placement constraints, and scaling mechanism.
     Link to deployment configuration for environment-owned sizing instead of copying replica or threshold values. -->

## Configuration

<!-- Link to the repository's authoritative configuration source when one exists.
     Describe notable settings and how environment overrides are applied.
     Do not assume every repository uses config/config.go or the same override path. -->

## Monitoring

<!-- Remove this section if not applicable. -->
<!-- Define exposed metrics and useful log searches.
     Describe key metrics and what an abnormal value means. -->

