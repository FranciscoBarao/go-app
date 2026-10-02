<!-- This template documents exactly one library. -->

# Library Name

<!-- Start the introduction with "Display Name (`package`)".
     Explain what the library provides and which problem it removes from its consumers in 3-5 sentences.
     State anything in the directory that is not importable code, such as a test helper or a generator.
     Keep a short library short. Do not pad a single-purpose helper to match the length of a larger one. -->

**Import Path:** <!-- The path consumers import -->

## Design

<!-- Describe how the library is organized: 
     If the library hides backends behind an interface, name the interface (facade) and the implementations (adapters).
     If it does not, say so in a sentence. Name the conceptual entry type and how sub-packages relate. 
     Document each sub-package in one or two lines and link to the code for depth, rather than giving each one its own README.
     Do not list constructors, option types, or sentinel errors here; those belong under Exports.
     A method-level description here goes stale on the next change, so do not embed a generated class diagram or
     per-package UML. -->

## Dependencies

<!-- Identify what this library requires: the infrastructure and external SDKs its implementations talk to, such as
     Redis, PostgreSQL, Pub/Sub, or object storage, and the sibling libraries in this repository that it builds on.
     Explain what each dependency provides and which sub-package uses it.
     Describe categories rather than enumerating module versions, which the dependency manifest owns. -->

## Consumers

<!-- Identify who imports this library and what they use it for.
     Name consumer roles rather than listing services/repositories.
     Don't maintain a full importer list.
     Record sibling libraries in this repository that depend on it, and any consumer outside the platform services -->

## Interfaces

<!-- For a library the exported surface is the interface. Record what is public, what a caller is expected to implement
     itself, and anything exported but not intended for general use. 
     Do not repeat the package layout from Design. -->

### Exports

<!-- Describe the consumer contract: constructors, option types, and sentinel errors.
     Name option types here; their fields and defaults belong under Configuration.
     Identify stability where it differs, including experimental and deprecated packages.
     Link to the implementation instead of reproducing signatures or a walkthrough; those belong under How to Use. -->


## Configuration

<!-- Remove this section if the library holds no configuration of its own. -->
<!-- Document only what this package defines: constructor parameters, option structs, and functional options, plus the
     defaults applied when a value is absent. Link to the type in this package.
     Do not document YAML, environment, or keys that a consuming service loads.
     Those belong in the service README.
     Do not copy credentials or environment-specific secret values into the README. -->

## How to Use

<!-- Numbered steps a consumer follows to use this library, in order: import, construct, supply required options,
     call the main operation, handle errors, and close or release anything that must be released.
     Show realistic paths, for every adapter or option. Keep code snippets short; a step that is obvious from the
     constructor does not need a snippet.
     Point at a test or example in this package when that is the maintained illustration, rather than duplicating it.
     Keep a short library to a few steps. -->

