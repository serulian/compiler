# Serulian Toolkit - Toolkit and Compiler for building Serulian projects

[![godoc](https://godoc.org/github.com/Serulian/compiler?status.svg)](http://godoc.org/github.com/Serulian/compiler)
[![Build Status (Travis)](https://travis-ci.org/Serulian/compiler.svg?branch=master)](https://travis-ci.org/Serulian/compiler)
[![Container Image on Quay](https://quay.io/repository/serulian/compiler/status "Container Image on Quay")](https://quay.io/repository/serulian/compiler)

The Serulian toolkit and compiler provides tooling for developing, building, formatting and testing [Serulian](https://github.com/Serulian/spec) web/mobile applications.

## Project Status

The Serulian toolkit is currently in **alpha** level development, which means it will be changing rapidly without guarentees of backwards compatibility. The toolkit itself is however fairly advanced in implementing the spec and various tooling, with significant testing of the compilation system already in place.

## Commands

### Building a project

To build a project, execute `build` with the entrypoint Serulian source file for that project:

```sh
./serulian build entrypointfile.seru
```

The project will be built and output as `entrypointfile.seru.js` and `entrypointfile.seru.js.map` in the current directory.

### Developing a project

To use the Serulian toolkit in an edit-refresh-compile development mode, run the `develop` command with the entrypoint Serulian source file for that project:

```sh
./serulian develop entrypointfile.seru
```

The toolkit will then start a webserver on the desired port (default `8080`):

```sh
Serving development server for project entrypointfile.seru on :8080 at /entrypointfile.seru.js
```

Add the following `<script>` tag to your application:

```html
<script type="text/javascript" src="http://localhost:8080/entrypointfile.seru.js"></script>
```

On page load (or refresh) the project will be recompiled, with compilation status and any errors or warnings displayed in the **web console**.

### Formatting source code

The Serulian toolkit command `format` can be used to reformat Serulian source code:

```
./serulian format somedir/...
```

The above will reformat all Serulian files found under the `somedir` directory.


### Working with imports

[Imports in Serulian](https://github.com/Serulian/spec/blob/master/proposals/ImportsAndPackages.md) are usually tied to a specific commit SHA or tagged version. The Serulian toolkit commands `freeze`, `unfreeze`, `update` and `upgrade` can be used to easily manage the versions of these imports.

#### Freeze

The `imports freeze` command can be used to rewrite the import to point to its current HEAD SHA:

```sh
./serulian imports freeze ./... github.com/Serulian/somelib
```

Contents of the matching source file after `freeze`:

```seru
from "github.com/Serulian/somelib:somesha" import SomeThing
```

#### Unfreeze

The `unfreeze` command can be used to rewrite imports back to HEAD, for real-time development:

```sh
./serulian imports unfreeze ./... github.com/Serulian/somelib
```

#### Update

If the library being imported is versioned using [Semantic Versioning](http://semver.org/), the additional command `update` can be used to update the import from an existing semantic version to a *minor* later version:

Given contents:

```seru
from "github.com/Serulian/somelib@v1.2.3" import SomeThing
```

Running:

```sh
./serulian imports update ./... github.com/Serulian/somelib
```

Contents of the matching source file after `update`:

```seru
from "github.com/Serulian/somelib@v1.3.0" import SomeThing
```

**Note:** The version will *not* be upgraded if the only change available is a major version change. To apply a major version change, use `upgrade`.

#### Upgrade

If the library being imported is versioned using [Semantic Versioning](http://semver.org/), the additional command `upgrade` can be used to upgrade the import to the latest *stable* version:

```sh
./serulian imports upgrade ./... github.com/Serulian/somelib
```

Contents of the matching source file after `upgrade`:

```seru
from "github.com/Serulian/somelib@v1.2.3" import SomeThing
```

### Testing a project

The Serulian toolkit can use one or more test runners to test Serulian projects. The current default runner is the [Karma test runner](https://karma-runner.github.io) with the [Jasmine](http://jasmine.github.io/).

#### Writing tests

Note: The following is subject to change in the near future.

Tests are Serulian source files ended with the suffix `_test.seru`. For example, a file `foo.seru` would have an associated test file named `foo_test.seru`.

All test files must have an entrypoint function named `TEST` that describes the various tests (using Jasmine test format) to be run. Note that **all tests must be asynchronous** (i.e. call the `done()` method when complete).

```seru
// Import the various Jasmine definitions. A jasmine.webidl defining these functions is required.
from webidl`jasmine` import describe
from webidl`jasmine` import it
from webidl`jasmine` import expect

/**
 * TEST defines the entrypoint function for describing all the tests.
 */
function<void> TEST() {
	// Describe a single test group.
	describe('Bool', function() {

		// Describe a single test.
		describe('equal', function() {

			// Add a requirement to be tested.
  		it(&'true should be equal to true', function(done function<void>()) {
   			expect(true).toBe(true);

   			// Mark the test's body as complete.
   			done()
	    })
    })
  })
}
```

#### Running tests

To run tests, execute the `test` command with the proper runner and entrypoint:

```sh
./serulian test karma ./...
```

The test runner plugin (in this case Karma) will ensure the necessary packages are installed and then run the specified tests.

## Running via container

A pre-built container image is always available. For example, the following with build a project via Docker. Note the mounting of the directory containing the project.

```sh
docker pull quay.io/serulian/compiler:latest
docker run -t -v /my/source/path:/project quay.io/serulian/compiler:latest build project/myfile.seru
```

