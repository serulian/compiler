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
this.Serulian = (function($global) {
  // Save the current script URL. This is used below when spawning web workers, as we need this
  // script URL in order to run itself.
  var $__currentScriptSrc = null;
  if (typeof $global.document === 'object') {
    $__currentScriptSrc = $global.document.currentScript.src;
  }

  // __serulian_internal defines methods used by the core library that require to work around
  // and with the type system.
  $global.__serulian_internal = {

    // autoNominalWrap automatically wraps (boxes) ES primitives into their associated
    // normal Serulian nominal types.
    'autoNominalWrap': function(k, v) {
      if (v == null) {
        return v;
      }
      var typeName = $t.toESType(v);
      switch (typeName) {
        case 'object':
          if (k != '') {
            return $t.nominalwrap(v, $a.mapping($t.any));
          }
          break;
        case 'array':
          return $t.nominalwrap(v, $a.slice($t.any));
        case 'boolean':
          return $t.nominalwrap(v, $a.bool);
        case 'string':
          return $t.nominalwrap(v, $a.string);
        case 'number':
          if (Math.ceil(v) == v) {
            return $t.nominalwrap(v, $a.int);
          }
          return $t.nominalwrap(v, $a.float64);
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
  var $it = function(name, index) {
      var tpe = new Function("return function " + name + "() {};")();
      tpe.$typeref = function() {
        return {
          'i': index
        };
      };

      return tpe;
  };

  // $t defines all helper literals and methods used under the type system.
  var $t = {
    // any special type.
    'any': $it('Any', 'any'),

    // void special type.
    'void': $it('Void', 'void'),

    // null special type.
    'null': $it('Null', 'null'),

    // toESType returns the ECMAScript type of the given object.
    'toESType': function(obj) {
      return ({}).toString.call(obj).match(/\s([a-zA-Z]+)/)[1].toLowerCase()
    },

    // cast performs a cast of the given value to the given type, throwing on
    // failure.
    'cast': function(value, type) {
      // TODO: implement cast checking.
      
      // Automatically box if necessary.
      if (type.$box) {
        return type.$box($t.unbox(value))
      }

      return value
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

    // box boxes the given raw data, wrapping it in the necessary structure
    // for the given type.
    'box': function(data, type) {
      if (type.$box) {
        return type.$box(data);
      }

      return data;
    },

    // unbox unboxes the given instance, returning the raw underlying data.
    'unbox': function(instance) {
      // TODO: make this efficient.
      return JSON.parse(JSON.stringify(instance));
    },

    // equals performs a comparison of the two values by calling the $equals operator on the
    // values, if any. Otherwise, uses a simple reference comparison.
    'equals': function(left, right, type) {
      if (left === right) {
        return $promise.resolve($t.box(true, $a['bool']));
      }

      if (left == null || right == null) {
        return $promise.resolve($t.box(false, $a['bool']));
      }

      // If we have a nominal wrapped native value, compare directly.
      if ($t.toESType(left) != 'object') {
        return $promise.resolve($t.box(left === right, $a['bool']));
      }

      if (type.$equals) {
        return type.$equals($t.box(left, type), $t.box(right, type));
      }

      return $promise.resolve($t.box(false, $a['bool']));
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
        var newInstance = Object.create(type.prototype);
        newInstance = type.apply(newInstance, arguments) || newInstance;
        return newInstance;
      };
    },

    // nominalroot returns the root object behind a nominal type. The returned instance is
    // guarenteed to not be a nominal type instance.
    'nominalroot': function(instance) {
      if (instance == null) {
          return null;
      }

      if (instance.hasOwnProperty('$wrapped')) {
        return $t.nominalroot(instance.$wrapped);
      }

      return instance;
    },

    // nominalwrap wraps an object with a nominal type.
    'nominalwrap': function(instance, type) {
      return type.new($t.nominalroot(instance))
    },

    // nominalwrap unwraps a nominal type one level. Unlike nominal root, the resulting
    // instance can be another nominal type.
    'nominalunwrap': function(instance) {
      return instance.$wrapped;
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
          var args = arguments;

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
        var args = Array.prototype.slice.call(arguments);

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
            var value = data['value'];
            var typeref = data['typeref'];
            if (typeref) {
              value = $t.box(value, $t.typeforref(typeref))
            }

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

    // property wraps a getter handler and optional setter handler into a single function
    // call.
    'property': function(getter, opt_setter) {
      var f = function() {
        if (arguments.length == 1) {
          return opt_setter.apply(this, arguments);
        } else {
          return getter.apply(this, arguments);
        }
      };

      f.$property = true;
      return f;
    },

    // dynamicaccess looks for the given name under the given object and returns it. If the
    // name was not found *OR* the object is null, returns null.
    'dynamicaccess': function(obj, name) {
      if (obj == null || obj[name] == null) {
        return null;
      }

      var value = obj[name];
      if (typeof value == 'function' && value.$property) {
        return $promise.wrap(function() {
          return value.apply(obj, arguments);
        });
      }

      return value
    },

    // nullcompare checks if the value is null and, if not, returns the value. Otherwise,
    // returns 'otherwise'.
    'nullcompare': function(value, otherwise) {
      return value == null ? otherwise : value;
    },

    // sm wraps a state machine handler function into a state machine object.
  	'sm': function(caller) {
  		return {
        resources: {},
  			current: 0,
  			next: caller,

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

  // $promise defines helper methods around constructing and managing ES promises.
  var $promise = {

    // build returns a Promise that invokes the given state machine.
  	'build': function(statemachine) {
  		return new Promise(function(resolve, reject) {
        statemachine.resolve = function(value) {
          statemachine.popall();
          statemachine.current = -1;
          resolve(value);
        };

        statemachine.reject = function(value) {
          statemachine.popall();
          statemachine.current = -1;
          reject(value);
        };

  			var continueFunc = function() {
  				if (statemachine.current < 0) {
  					return;
  				}

  				statemachine.next(callFunc);				
			  };

        var callFunc = function() {
    			continueFunc();
          if (statemachine.current < 0) {
            statemachine.resolve(null);
          }
        };

        callFunc();
  		});
  	},

  	'all': function(promises) {
  		return Promise.all(promises);
  	},

    'new': function(f) {
      return new Promise(f);
    },

  	'empty': function() {
  		return new Promise(function() {
  			resolve();
  		});
  	},

    'resolve': function(value) {
      return Promise.resolve(value);
    },

  	'wrap': function(func) {
  		return Promise.resolve(func());
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

  // moduleInits defines a collection of all promises to initialize the various modules.
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
      return function(name, hasGenerics, alias, creator) {
        var buildType = function(n, args) {
          var args = args || [];

          // Create the type function itself, with the type's name.
          var tpe = new Function("return function " + n + "() {};")();

          // Add a way to retrieve a type ref for the type.
          tpe.$typeref = function() {
            var generics = [];
            for (var i = 0; i < args.length; ++i) {
              generics.push(args[i].$typeref())
            }

            return {
              't': moduleName + '.' + name,
              'g': generics
            };
          };

          // Build the type's static and prototype.
          creator.apply(tpe, args);

          // Add default type-system members.
          if (kind == 'type') {
            // toJSON.
            tpe.prototype.toJSON = function() {
              var root = $t.nominalroot(this);
              if (root.toJSON) {
                return root.toJSON()
              }

              return root;
            };
          } else if (kind == 'struct') {
            // $box.
            tpe.$box = function(data, opt_lazycheck) {
              var instance = new tpe();
              instance.$data = data;
              instance.$lazycheck = opt_lazycheck || true;
              return instance;
            };

            // toJSON.
            tpe.prototype.toJSON = function() {
              return this.$data;
            };

            // String.
            tpe.prototype.String = function() {
              return $promise.resolve($t.nominalwrap(JSON.stringify(this, null, ' '), $a['string']));
            };

            // Stringify.
            tpe.prototype.Stringify = function(T) {
              var $this = this;
              return function() {
                // Special case JSON, as it uses an internal method.
                if (T == $a['json']) {
                  return $promise.resolve($t.nominalwrap(JSON.stringify($this), $a['string']));
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
                  var parsed = JSON.parse($t.nominalunwrap(value));
                  var boxed = tpe.$box(parsed, true);

                  // Call Mapping to ensure every field is checked.
                  return boxed.Mapping().then(function() {
                    boxed.$lazycheck = false;
                    return $promise.resolve(boxed);
                  });
                }

                return T.Get().then(function(resolved) {
                  return (resolved.Parse(value)).then(function(parsed) {
                    // TODO: *efficiently* unwrap internal nominal types.
                    var data = JSON.parse(JSON.stringify($t.nominalunwrap(parsed)));
                    return $promise.resolve(tpe.$box(data));
                  });
                });
              };
            };
          }

          return tpe;
        };

        // Define the type on the module.
        if (hasGenerics) {
          module[name] = function(__genericargs) {
            var fullName = name;
            for (var i = 0; i < arguments.length; ++i) {
              fullName = fullName + '_' + arguments[i].name;
            }

            var tpe = buildType(name, arguments);
            tpe.$generic = arguments.callee;
            return tpe;
          };
        } else {
          module[name] = buildType(name);
        }

        // If the type has an alias, add it to the global alias map.
        if (alias) {
          $a[alias] = module[name];
        }
      };
    };

    // $init adds a promise to the module inits array.
    module.$init = function(cpromise) {
      moduleInits.push(cpromise);
    };

    module.$struct = $newtypebuilder('struct');
  	module.$class = $newtypebuilder('class');
  	module.$interface = $newtypebuilder('interface');
    module.$type = $newtypebuilder('type');

  	creator.call(module)
  };

  {{ range $idx, $kv := .Iter }}
  	{{ $kv.Value }}
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
          var arguments = data['arguments'];
          var method = $w[methodId];

          var send = function(kind) {
            return function(value) {
              var message = {
                'token': token,
                'value': value,
                'kind': kind
              };

              // If the object is a Serulian object, then add its type ref.
              if (value != null && value.constructor.$typeref) {
                message['typeref'] = value.constructor.$typeref();
                message['value'] = $t.unbox(value);
              }

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

          method.apply(null, arguments).then(send('resolve')).catch(send('reject'));
          break;
      }
    };
  };

  // Return a promise which initializes all modules and, once complete, returns the global
  // namespace map.
  return $promise.all(moduleInits).then(function() {
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
