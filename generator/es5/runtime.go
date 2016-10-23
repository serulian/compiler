// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

// Note: nativenew is based on http://www.bennadel.com/blog/2291-invoking-a-native-javascript-constructor-using-call-or-apply.htm
// Note: toESType is based on https://javascriptweblog.wordpress.com/2011/08/08/fixing-the-javascript-typeof-operator/
// Note: uuid generation from https://stackoverflow.com/questions/105034/create-guid-uuid-in-javascript

// runtimeTemplate contains all the necessary code for wrapping the generated modules into a complete Serulian
// runtime bundle.
const runtimeTemplate = `
"use strict";
this.Serulian = (function($global) {
  var BOXED_DATA_PROPERTY = '$wrapped';

  // Save the current script URL. This is used below when spawning web workers, as we need this
  // script URL in order to run itself.
  var $__currentScriptSrc = null;
  if (typeof $global.document === 'object') {
    $__currentScriptSrc = $global.document.currentScript.src;
  }

  // __serulian_internal defines methods used by the core library that require to work around
  // and with the type system.
  $global.__serulian_internal = {
    // autoUnbox automatically unboxes ES primitives into their raw data.
    'autoUnbox': function(k, v) {
      return $t.unbox(v);
    },

    // autoBox automatically boxes ES primitives into their associated normal Serulian nominal types.
    'autoBox': function(k, v) {
      if (v == null) {
        return v;
      }
      var typeName = $t.toESType(v);
      switch (typeName) {
        case 'object':
          if (k != '') {
            return $t.fastbox(v, $a.mapping($t.any));
          }
          break;

        case 'array':
          return $t.fastbox(v, $a.slice($t.any));

        case 'boolean':
          return $t.fastbox(v, $a.bool);

        case 'string':
          return $t.fastbox(v, $a.string);

        case 'number':
          if (Math.ceil(v) == v) {
            return $t.fastbox(v, $a.int);
          }
          return $t.fastbox(v, $a.float64);
      }
      return v;
    }
  };

  // $g is defines the root of the type paths in Serulian. All modules will be placed somewhere
  // in a tree starting at $g.
  var $g = {};

  // $a defines a map of aliases to their types.
  var $a = {};

  // $w defines a map from unique web-worker-function UUIDs to the corresponding function. Used
  // for lookup by the web workers when async functions are executed across the wire.
  var $w = {};

  // $it defines an internal type with a name. An internal type is any type that doesn't have
  // real implementation (such as 'any', 'void', 'null', etc).
  var $it = function(name, typeIndex) {
      var tpe = new Function("return function " + name + "() {};")();
      tpe.$typeId = typeIndex;
      tpe.$typeref = function() {
        return {
          'i': typeIndex
        };
      };

      return tpe;
  };

  // $t defines all helper literals and methods used under the type system.
  var $t = {
    // any special type.
    'any': $it('Any', 'any'),

    // struct special type.
    'struct': $it('Struct', 'struct'),

    // void special type.
    'void': $it('Void', 'void'),

    // null special type.
    'null': $it('Null', 'null'),

    // toESType returns the ECMAScript type of the given object.
    'toESType': function(obj) {
      return ({}).toString.call(obj).match(/\s([a-zA-Z]+)/)[1].toLowerCase()
    },

    // From: http://stackoverflow.com/a/15714445
    'functionName': function(func) {
      if (func.name) {
        return func.name;
      }

      var ret = func.toString();
      ret = ret.substr('function '.length);
      ret = ret.substr(0, ret.indexOf('('));
      return ret;
    },

    // typeid returns the globally unique ID for the given type.
    'typeid': function(type) {
      return type.$typeId || $t.functionName(type);
    },

    // buildDataForValue builds an object containing the given value in its unboxed form,
    // as well as its Serulian type reference information (if any).
    'buildDataForValue': function(value) {
      if (value == null) {
        return {
          'v': null
        }
      }

      if (value.constructor.$typeref) {
        return {
          'v': $t.unbox(value),
          't': value.constructor.$typeref()
        }
      } else {
        return {
          'v': value
        }
      }
    },

    // buildValueFromData builds a Serulian value from its unboxed form and optional
    // type reference information.
    'buildValueFromData': function(data) {
      if (!data['t']) {
        return data['v'];
      }

      return $t.box(data['v'], $t.typeforref(data['t']));
    },

    // unbox returns the root object behind a nominal or structural instance.
    'unbox': function(instance) {
      if (instance != null && instance.hasOwnProperty(BOXED_DATA_PROPERTY)) {
        return instance[BOXED_DATA_PROPERTY];
      }

      return instance;
    },

    // box wraps an object with a nominal or structural type.
    'box': function(instance, type) {
      if (instance == null) {
        return null;
      }
      
      return type.$box($t.unbox(instance));
    },

    // fastbox wraps an object with a nominal or structural type, without first
    // unboxing or null checking. This method will fail if those two assumptions
    // are not correct.
    'fastbox': function(instance, type) {
      return type.$box(instance);
    },

    // roottype returns the root type of the given type. If the type is a nominal
    // type, then its root is returned. Otherwise, the type itself is returned.
    'roottype': function(type) {
      if (type.$roottype) {
        return type.$roottype();
      }

      return type;
    },

    // istype returns true if the specified value can be used in place of the given type.
    'istype': function(value, type) {
      // Quick check to see if it matches directly or is 'any'.
      if (type == $t.any || (value != null && (value.constructor == type || value instanceof type))) {
        return true;
      }

      // Quick check for struct.
      if (type == $t.struct) {
        // Find the root type of the current value's type. If it is structural or a serializable
        // native type, then we can cast it to struct.
        var roottype = $t.roottype(value.constructor);
        return (roottype.$typekind == 'struct' || 
                roottype == Number ||
                roottype == String ||
                roottype == Boolean ||
                roottype == Object);
      }

      // Quick check for function.
      // TODO: this needs to properly check the function's parameter and return types
      // once that information is added to the value. For now, we just make sure we have
      // a function.
      if (type.$generic == $a['function'] && typeof value == 'function') {
        return value;
      }

      var targetKind = type.$typekind;
      switch (targetKind) {
        case 'struct':
        case 'type':
        case 'class':
          // Direct matching is the requirement.
          return false;

        case 'interface':
          // Check if the other type implements the interface by comparing type signatures.
          var targetSignature = type.$typesig();
          var valueSignature = value.constructor.$typesig();

          var expectedKeys = Object.keys(targetSignature);
          for (var i = 0; i < expectedKeys.length; ++i) {
            var expectedKey = expectedKeys[i];
            if (valueSignature[expectedKey] !== true) {
              return false;
            }
          }
          return true;

        default:
          return false;
      }
    },

    // cast performs a cast of the given value to the given type, throwing on
    // failure.
    'cast': function(value, type, opt_allownull) {
      // Ensure that we don't cast to non-nullable if null.
      if (value == null && !opt_allownull) {
        throw Error('Cannot cast null value to ' + type.toString())        
      }

      // Check if the type matches directly.
      if ($t.istype(value, type)) {
        return value;
      }

      // Otherwise handle some cast-only cases and errors.
      var targetKind = type.$typekind;
      switch (targetKind) {
        case 'struct':
          // Note: This cast can result in us getting a struct over invalid data, but
          // the struct will ensure it is correct anyway, so this is allowed.
          if (value.constructor == Object) {
            break;
          }

          throw Error('Cannot cast ' + value.constructor.toString() + ' to ' + type.toString());

        case 'class':
        case 'interface':
          // Since we check for equality above, this must fail.
          throw Error('Cannot cast ' + value.constructor.toString() + ' to ' + type.toString());

        case 'type':
          // Casting only is allowed if the nominal and the existing value's
          // type have the same root type.
          if ($t.roottype(value.constructor) != $t.roottype(type)) {
            throw Error('Cannot auto-box ' + value.constructor.toString() + ' to ' + type.toString());
          }
          break;

        case undefined:
          // Cannot cast a non-native type to a native type. The native-to-native check occurs in
          // istype above.
          throw Error((('Cannot cast ' + value.constructor.toString()) + ' to ') + type.toString());
      }

      // Automatically box if necessary.
      if (type.$box) {
        return $t.box(value, type);
      }

      return value
    },

    // equals performs a comparison of the two values by calling the $equals operator on the
    // values, if any. Otherwise, uses a simple reference comparison.
    'equals': function(left, right, type) {
      // Check for direct equality.
      if (left === right) {
        return $promise.resolve($t.fastbox(true, $a['bool']));
      }

      // Check for null.
      if (left == null || right == null) {
        return $promise.resolve($t.fastbox(false, $a['bool']));
      }

      // Check for defined equals operator.
      if (type.$equals) {
        return type.$equals($t.box(left, type), $t.box(right, type));
      }

      // Otherwise we cannot compare, so we treat the objects as not equal.
      return $promise.resolve($t.fastbox(false, $a['bool']));
    },

    // ensurevalue ensures that the given value is of the given type. If not,
    // raises an exception.
    'ensurevalue': function(value, type, canBeNull, name) {
      if (value == null) {
        if (!canBeNull) {
          throw Error('Missing value for non-nullable field ' + name)
        }
        return;
      }

      var check = function(serutype, estype) {
        if (type == $a[serutype] || type.$generic == $a[serutype]) {
          if ($t.toESType(value) != estype) {
            throw Error('Expected ' + serutype + ' for field ' + name + ', found: ' + $t.toESType(value))
          }
          return true;
        }

        return false;
      };

      if (check('string', 'string')) { return; }
      if (check('float64', 'number')) { return; }
      if (check('int', 'number')) { return; }
      if (check('bool', 'boolean')) { return; }
      if (check('slice', 'array')) { return; }

      if ($t.toESType(value) != 'object') {
        throw Error('Expected object for field ' + name + ', found: ' + $t.toESType(value))        
      }
    },

    // nativenew creates a new instance of the *ECMAScript* type specified (e.g. Number, String).
    'nativenew': function(type) {
      return function () {
        if (arguments.length == 0) {
          return new type();
        }

        if (type == $global.Promise) {
          return new Promise(arguments[0]);
        }

        var newInstance = Object.create(type.prototype);
        newInstance = type.apply(newInstance, arguments) || newInstance;
        return newInstance;
      };
    },

    // typeforref deserializes a typeref into a local type.
    'typeforref': function(typeref) {
      if (typeref['i']) {
        return $t[typeref['i']];
      }

      // Lookup the type.
      var parts = typeref['t'].split('.');
      var current = $g;
      for (var i = 0; i < parts.length; ++i) {
        current = current[parts[i]];
      }

      if (!typeref['g'].length) {
        return current;
      }

      // Get generics.
      var generics = typeref['g'].map(function(generic) {
        return $t.typeforref(generic);
      });

      // Apply the generics to the type.
      return current.apply(current, generics);
    },

    // uuid returns a new *cryptographically secure* UUID.
    'uuid': function() {
        var buf = new Uint16Array(8);
        crypto.getRandomValues(buf);
        var S4 = function(num) {
            var ret = num.toString(16);
            while(ret.length < 4){
                ret = "0"+ret;
            }
            return ret;
        };
        return (S4(buf[0])+S4(buf[1])+"-"+S4(buf[2])+"-"+S4(buf[3])+"-"+S4(buf[4])+"-"+S4(buf[5])+S4(buf[6])+S4(buf[7]));
    },

    // defineStructField defines a new field on a structural type.
    'defineStructField': function(structType, name, serializableName, typeref, opt_nominalRootType, opt_nullAllowed) {
      var field = {
        'name': name,
        'serializableName': serializableName,
        'typeref': typeref,
        'nominalRootTyperef': opt_nominalRootType || typeref,
        'nullAllowed': opt_nullAllowed
      };
      structType.$fields.push(field);

      Object.defineProperty(structType.prototype, name, {
        get: function() {
          // If the underlying object was not created by the runtime, then we must typecheck
          // to be safe when accessing the field.
          var boxedData = this[BOXED_DATA_PROPERTY];
          if (!boxedData.$runtimecreated) {
            if (!this.$lazychecked[field.name]) {
              $t.ensurevalue($t.unbox(boxedData[field.serializableName]), field.nominalRootTyperef(), field.nullAllowed, field.name);
              this.$lazychecked[field.name] = true;
            }

            var fieldType = field.typeref();
            if (fieldType.$box) {
              return $t.box(boxedData[field.serializableName], fieldType);
            } else {
              return boxedData[field.serializableName];
            }
          }

          // Otherwise, simply return the field.
          return boxedData[name];
        },

        set: function(value) {
          this[BOXED_DATA_PROPERTY][name] = value;
        }
      });
    },

    // workerwrap wraps a function definition to be executed via a web worker. When the function
    // is invoked a new web worker will be spawned on this script. The method ID and all arguments
    // will be serialized and sent to the web worker, which will lookup the function, invoke it,
    // await the promise result and send the data back to this script.
    'workerwrap': function(methodId, f) {
      // Save the method by its unique ID.
      $w[methodId] = f;

      // If already inside a worker, return a function to execute asynchronously locally.
      // TODO: Chrome does not support nested workers but Gecko does, so it might be worth
      // feature checking here.
      if (!$__currentScriptSrc) {
        return function() {
          var $this = this;

          // Reference: https://github.com/petkaantonov/bluebird/wiki/Optimization-killers#3-managing-arguments
          var args = new Array(arguments.length);
          for (var i = 0; i < args.length; ++i) {
              args[i] = arguments[i];
          }

          var promise = new Promise(function(resolve, reject) {
            $global.setTimeout(function() {
              f.apply($this, args).then(function(value) {
                resolve(value);
              }).catch(function(value) {
                reject(value);
              });
            }, 0);
          });
          return promise;
        };
      }

      // Otherwise return a function to execute via a worker.
      return function() {
        var token = $t.uuid();
        var args = Array.prototype.map.call(arguments, $t.buildDataForValue);

        var promise = new Promise(function(resolve, reject) {
          // Start a new worker on the current script with the token.
          var worker = new Worker($__currentScriptSrc + "?__serulian_async_token=" + token);

          worker.onmessage = function(e) {
            // Ensure we received a trusted message.
            if (!e.isTrusted) {
              worker.terminate();
              return;
            }

            // Ensure that this is the result for the sent token.
            var data = e.data;
            if (data['token'] != token) {
              return;
            }

            // Box the value if necessary.
            var value = $t.buildValueFromData(data['value']);

            // Report the result.
            var kind = data['kind'];
            if (kind == 'resolve') {
              resolve(value);
            } else {
              reject(value);
            }

            // Terminate the worker in case it has not yet closed.
            worker.terminate();
          };

          // Tell the worker to invoke the function, with its token and arguments.
          worker.postMessage({
            'action': 'invoke',
            'arguments': args,
            'method': methodId,
            'token': token
          });           
         });
         return promise;
      };
    },

    // property wraps a getter handler and marks it as belonging to a property.
    'property': function(getter) {
      getter.$property = true;
      return getter;
    },

    // nullableinvoke invokes the function found at the given name on the object, but
    // only if the object is not null.
    'nullableinvoke': function(obj, name, promising, args) {
      var found = obj != null ? obj[name] : null;
      if (found == null) {
        return promising ? $promise.resolve(null) : null;
      }

      return found.apply(obj, args);
    },

    // dynamicaccess looks for the given name under the given object and returns it, wrapped in
    // a promise (if not already a promise). If the name was not found *OR* the object is null,
    // returns a promise resolving null.
    'dynamicaccess': function(obj, name) {
      if (obj == null || obj[name] == null) {
        return $promise.resolve(null);
      }

      var value = obj[name];
      if (typeof value == 'function') {
        if (value.$property) {
          return value.apply(obj, arguments);
        } else {
          return $promise.resolve(function() {
            return value.apply(obj, arguments);
          });
        }
      }

      return $promise.resolve(value);
    },

    // assertnotnull checks if the value is null and, if so, raises an error. Otherwise,
    // returns the value.
    'assertnotnull': function(value) {
      if (value == null) {
        throw Error('Value should not be null')
      }

      return value;
    },

    // nullcompare checks if the value is null and, if not, returns the value. Otherwise,
    // returns 'otherwise'.
    'nullcompare': function(value, otherwise) {
      return value == null ? otherwise : value;
    },

    // resourcehandler returns a function for handling resources in a function.
  	'resourcehandler': function() {
  		return {
        resources: {},
        bind: function(func) {
          if (func.$__resourcebound) {
            return func;
          }

          var r = this;
          var f = function() {
            r.popall();
            return func.apply(this, arguments);
          };

          f.$__resourcebound = true;
          return f;
        },

        pushr: function(value, name) {
          this.resources[name] = value;
        },

        popr: function(names) {
          var promises = [];

          for (var i = 0; i < arguments.length; ++i) {
            var name = arguments[i];
            if (this.resources[name]) {
              promises.push(this.resources[name].Release());
              delete this.resources[name];
            }
          }

          if (promises.length > 0) {
            return $promise.all(promises);
          } else {
            return $promise.resolve(null);
          }
        },

        popall: function() {
          for (var name in this.resources) {
            if (this.resources.hasOwnProperty(name)) {
              this.resources[name].Release();
            }
          }
        }
  		};
  	}
  };

  // $generator defines helper methods around constructing generators for Streams.
  var $generator = {
    // directempty returns a new empty generator.
    'directempty': function() {
      var stream = {
        'Next': function() {
           return $promise.new(function (resolve, reject) {
             $a['tuple']($t.any, $a['bool']).Build(null, false).then(resolve);
           });
        },
      };
      return stream;
    },

    // empty returns a new empty stream via a promise.
    'empty': function() {
      return $promise.resolve($generator.directempty());
    },

    // new returns a stream wrapping the given generator function.
    'new': function (f) {
      var stream = {
        '$is': null,
        'Next': function () {
          return $promise.new(function (resolve, reject) {
            if (stream.$is != null) {
              stream.$is.Next().then(function (tuple) {
                if ($t.unbox(tuple.Second)) {
                  resolve(tuple);
                } else {
                  stream.$is = null;
                  stream.Next().then(resolve, reject);
                }
              }).catch(function (rejected) {
                reject(rejected);
              });
              return;
            }

            var $yield = function (value) {
              $a['tuple']($t.any, $a['bool']).Build(value, $t.fastbox(true, $a['bool'])).then(resolve);
            };
            
            var $done = function () {
              $a['tuple']($t.any, $a['bool']).Build(null, $t.fastbox(false, $a['bool'])).then(resolve);
            };

            var $yieldin = function (ins) {
              stream.$is = ins;
              stream.Next().then(resolve, reject);
            };
            
            f($yield, $yieldin, reject, $done);
          });
        },
      };
      return $promise.resolve(stream);
    }
  };

  // $promise defines helper methods around constructing and managing ES promises.
  var $promise = {
  	'all': function(promises) {
  		return Promise.all(promises);
  	},

    'new': function(f) {
      return new Promise(f);
    },

  	'empty': function() {
  		return Promise.resolve(null);
  	},

    'resolve': function(value) {
      return Promise.resolve(value);
    },

  	'wrap': function(func) {
  		return Promise.resolve(func());
  	},

    // shortcircuit returns a promise that resolves the given boolean value if and only if
    // it is not equal to the right value. Returns null otherwise.
    'shortcircuit': function(left, right) {
      if (left != right) {
        return $promise.resolve(left);
      }
    },

    // translate translates a Serulian Promise into an ES promise.
    'translate': function(prom) {
       if (!prom.Then) {
         return prom;
       }

       return {
          'then': function() {
             return prom.Then.apply(prom, arguments);
          },
          'catch': function() {
             return prom.Catch.apply(prom, arguments);
          }
       };
    }
  };

  // moduleInits defines a collection of functions that, when called, return promises to initialize the
  // various modules.
  var moduleInits = [];

  // $module defines a module in the type system.
  var $module = function(moduleName, creator) {
  	var module = {};

    // Define the module under the gloal path array.
    var parts = moduleName.split('.');
    var current = $g;
    for (var i = 0; i < parts.length - 1; ++i) {
      if (!current[parts[i]]) {
        current[parts[i]] = {};
      }
      current = current[parts[i]]
    }

    current[parts[parts.length - 1]] = module;

    // $newtypebuilder is a helper function for creating types of a particular kind. Returns
    // a function that can be used to create a type of the specified kind.
    var $newtypebuilder = function(kind) {
      return function(typeId, name, hasGenerics, alias, creator) {
        var buildType = function(fullTypeId, fullName, args) {
          var args = args || [];

          // Create the type function itself, with the type's name.
          var tpe = new Function("return function " + fullName + "() {};")();

          // Add a way to retrieve a type ref for the type.
          tpe.$typeref = function() {
            if (!hasGenerics) {
              return {
                't': moduleName + '.' + name
              };
            }

            var generics = [];
            for (var i = 0; i < args.length; ++i) {
              generics.push(args[i].$typeref())
            }

            return {
              't': moduleName + '.' + name,
              'g': generics
            };
          };

          tpe.$typeId = fullTypeId;
          tpe.$typekind = kind;

          // Build the type's static and prototype.
          creator.apply(tpe, args);

          // Add default type-system members.
          if (kind == 'struct') {
            // $box.
            tpe.$box = function(data) {    
              var instance = new tpe();
              instance[BOXED_DATA_PROPERTY] = data;
              instance.$lazychecked = {};
              return instance;
            };

            // $markruntimecreated marks the struct as having been created by the runtime, allowing
            // for certain optimizations to be used and checks to be turned off.
            tpe.prototype.$markruntimecreated = function() {
              Object.defineProperty(this[BOXED_DATA_PROPERTY], '$runtimecreated', {
                enumerable: false,
                configurable: true,
                value: true
              });
            };

            // String.
            tpe.prototype.String = function() {
              return $promise.resolve($t.fastbox(JSON.stringify(this, $global.__serulian_internal.autoUnbox, ' '), $a['string']));
            };

            // Clone.
            tpe.prototype.Clone = function() {
              var instance = new tpe();
              if (Object.assign) {
                instance[BOXED_DATA_PROPERTY] = Object.assign({}, this[BOXED_DATA_PROPERTY]);
              } else {
                instance[BOXED_DATA_PROPERTY] = {};
                for (var key in this[BOXED_DATA_PROPERTY]) {
                  if (this[BOXED_DATA_PROPERTY].hasOwnProperty(key)) {
                    instance[BOXED_DATA_PROPERTY][key] = this[BOXED_DATA_PROPERTY][key];
                  }
                }
              }

              if (this[BOXED_DATA_PROPERTY].$runtimecreated) {
                instance.$markruntimecreated();
              }

              return $promise.resolve(instance);
            };

            // Stringify.
            tpe.prototype.Stringify = function(T) {
              var $this = this;
              return function() {
                // Special case JSON, as it uses an internal method.
                if (T == $a['json']) {
                  return $promise.resolve($t.fastbox(JSON.stringify($this, $global.__serulian_internal.autoUnbox), $a['string']));
                }

                return $this.Mapping().then(function(mapped) {
                  return T.Get().then(function(resolved) {
                    return resolved.Stringify(mapped);
                  });
                });
              };
            };

            // Parse.
            tpe.Parse = function(T) {
              return function(value) {
                // Special case JSON for performance, as it uses an internal method.
                if (T == $a['json']) {
                  var parsed = JSON.parse($t.unbox(value));
                  var boxed = $t.fastbox(parsed, tpe);

                  // Call Mapping to ensure every field is checked.
                  return boxed.Mapping().then(function() {
                    return $promise.resolve(boxed);
                  });
                }

                return T.Get().then(function(resolved) {
                  return (resolved.Parse(value)).then(function(parsed) {
                    return $promise.resolve($t.box(parsed, tpe));
                  });
                });
              };
            };

            // Equals.
            tpe.$equals = function(left, right) {
              if (left === right) {
                return $promise.resolve($t.fastbox(true, $a['bool']));
              }

              // TODO: find a way to do this without checking *all* fields.
              var promises = [];
              tpe.$fields.forEach(function(field) {
                promises.push(
                  $t.equals(left[BOXED_DATA_PROPERTY][field.serializableName], 
                            right[BOXED_DATA_PROPERTY][field.serializableName],
                            field.typeref()));
              });

              return Promise.all(promises).then(function(values) {
                for (var i = 0; i < values.length; i++) {
                  if (!$t.unbox(values[i])) {
                    return values[i];
                  }
                }

                return $t.fastbox(true, $a['bool']);
              });
            };

            // Mapping.
            tpe.prototype.Mapping = function() {
              if (this.$serucreated) {
                // Fast-path for compiler-constructed instances. All data is guarenteed to already
                // be boxed.
                return $promise.resolve($t.fastbox(this[BOXED_DATA_PROPERTY], $a['mapping']($t.any)));
              } else {
                // Slower path for instances unboxed from native data. We call the properties
                // to make sure we have the boxed forms.
                var $this = this;
                var mapped = {};
                tpe.$fields.forEach(function(field) {
                  mapped[field.serializableName] = $this[field.name];
                });

                return $promise.resolve($t.fastbox(mapped, $a['mapping']($t.any)));
              }
            };
          } // end struct

          return tpe;
        };

        // Define the type on the module.
        if (hasGenerics) {
          module[name] = function genericType() {
            var fullName = name;
            var fullId = typeId;

            // Reference: https://github.com/petkaantonov/bluebird/wiki/Optimization-killers#3-managing-arguments
            var generics = new Array(arguments.length);
            for (var i = 0; i < generics.length; ++i) {
              fullName = fullName + '_' + $t.functionName(arguments[i]);
              if (i == 0) {
                fullId = fullId + '<';
              } else {
                fullId = fullId + ',';
              }

              fullId = fullId + arguments[i].$typeId;
              generics[i] = arguments[i];
            }

            // Check for a cached version of the generic type.
            var cached = module[fullName];
            if (cached) {
              return cached;
            }

            var tpe = buildType(fullId + '>', fullName, generics);
            tpe.$generic = genericType;
            return module[fullName] = tpe;
          };
        } else {
          module[name] = buildType(typeId, name);
        }

        // If the type has an alias, add it to the global alias map.
        if (alias) {
          $a[alias] = module[name];
        }
      };
    };

    // $init adds a promise to the module inits array.
    module.$init = function(callback, fieldId, dependencyIds) {
      moduleInits.push({
        'callback': callback,
        'id': fieldId,
        'depends': dependencyIds
      });
    };

    module.$struct = $newtypebuilder('struct');
  	module.$class = $newtypebuilder('class');
  	module.$interface = $newtypebuilder('interface');
    module.$type = $newtypebuilder('type');

  	creator.call(module)
  };

  {{ range $idx, $kv := .UnsafeIter }}
  	{{ emit $kv.Value }}
  {{ end }}

  // $executeWorkerMethod executes an async called function in this web worker. When invoked with
  // a call token, the method will add an onmessage listener, receive the message with the function
  // to invoke, invoke the function, and send the result back to the caller, closing the web worker
  // once complete.
  $g.$executeWorkerMethod = function(token) {
    $global.onmessage = function(e) {
      // Ensure we have a trusted message.
      if (!e.isTrusted) {
        $global.close();
        return;
      }

      // Ensure this web worker is for the expected token.
      var data = e.data;
      if (data['token'] != token) {
        throw Error('Invalid token')
        $global.close();
      }

      switch (data['action']) {
        case 'invoke':
          var methodId = data['method'];
          var method = $w[methodId];

          var args = data['arguments'].map($t.buildValueFromData);
          var send = function(kind) {
            return function(value) {
              var message = {
                'token': token,
                'value': $t.buildDataForValue(value),
                'kind': kind
              };

              // Try to send the message to the other process. If it fails, then
              // the rejected value is most likely some sort of local exception,
              // so we just throw it.
              try {
                $global.postMessage(message);
              } catch (e) {
                if (kind == 'reject') {
                  throw value;
                } else {
                  // Should never happen, but just in case.
                  throw e;
                }
              }

              $global.close();
            };
          };

          method.apply(null, args).then(send('resolve')).catch(send('reject'));
          break;
      }
    };
  };

  var buildPromises = function(items) {
    var seen = {};
    var result = [];

    // Build a map by item ID.
    var itemsById = {};
    items.forEach(function(item) {
      itemsById[item.id] = item;
    });

    // Topo-sort and execute in order.
    items.forEach(function visit(item) {
      if (seen[item.id]) {
        return;
      }

      seen[item.id] = true;
      item.depends.forEach(function(depId) {
        visit(itemsById[depId]);
      });
      
      item['promise'] = item['callback']();
    });

    // Build the dep-chained promises.
    return items.map(function(item) {
      if (!item.depends.length) {
        return item['promise'];
      }

      var current = $promise.resolve();
      item.depends.forEach(function(depId) {
        current = current.then(function(resolved) {
          return itemsById[depId]['promise'];
        });
      });

      return current.then(function(resolved) {
        return item['promise'];
      });
    });
  };

  // Return a promise which initializes all modules and, once complete, returns the global
  // namespace map.
  return $promise.all(buildPromises(moduleInits)).then(function() {
  	return $g;
  });
})(this)

// Handle web-worker calls.
if (typeof importScripts === 'function') {
  var runWorker = function() {
    var search = location.search;
    if (!search || search[0] != '?') {
      return;
    }

    var searchPairs = search.substr(1).split('&')
    if (searchPairs.length < 1) {
      return;
    }

    for (var i = 0; i < searchPairs.length; ++i) {
      var pair = searchPairs[i].split('=');
      if (pair[0] == '__serulian_async_token') {
        this.Serulian.then(function(global) {
          global.$executeWorkerMethod(pair[1]);
        });
        return;
      }
    }

    close();
  };
  runWorker();
}
`
