// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

// Note: nativenew is based on http://www.bennadel.com/blog/2291-invoking-a-native-javascript-constructor-using-call-or-apply.htm
// Note: toType is based on https://javascriptweblog.wordpress.com/2011/08/08/fixing-the-javascript-typeof-operator/
// Note: uuid generation from https://stackoverflow.com/questions/105034/create-guid-uuid-in-javascript

// runtimeTemplate contains all the necessary code for wrapping generated modules into a complete Serulian
// runtime bundle.
const runtimeTemplate = `
this.Serulian = (function($global) {
  var $__currentScriptSrc = '';
  if (typeof document === 'object') {
    $__currentScriptSrc = document.currentScript.src;
  } else {
    // TODO: fix for nested calls to workers by workers.
  }

  $global.__serulian_internal = {
    'autoNominalWrap': function(k, v) {
      if (v == null) {
        return v;
      }

      var typeName = $t.toType(v);
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

  var $g = {};
  var $a = {};
  var $w = {};

  var $t = {
    'toType': function(obj) {
      return ({}).toString.call(obj).match(/\s([a-zA-Z]+)/)[1].toLowerCase()
    },

    'any': new Function("return function any() {};")(),

    'cast': function(value, type) {
      // TODO: implement cast checking.
      return value
    },

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

    'nativenew': function(type) {
      return function () {
        var newInstance = Object.create(type.prototype);
        newInstance = type.apply(newInstance, arguments) || newInstance;
        return newInstance;
      };
    },

    'nominalroot': function(instance) {
      if (instance.$wrapped) {
        return instance.$wrapped;
      }

      return instance;
    },

    'nominalwrap': function(instance, type) {
      return type.new(instance)
    },

    'nominalunwrap': function(instance) {
      return instance.$wrapped;
    },

    'workerwrap': function(methodId, f) {
      $w[methodId] = f;
      return function() {
         var args = Array.prototype.slice.call(arguments);
         var token = $t.uuid();

         var promise = new Promise(function(resolve, reject) {
           var worker = new Worker($__currentScriptSrc + "?__serulian_async_token=" + token);
           worker.onmessage = function(e) {
              if (!e.isTrusted) { return; }

              var data = e.data;
              if (data['token'] != token) {
                return;
              }

              if (data['result']) {
                resolve(data['result']);
              } else {
                reject(data['reject']);
              }
           };

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

    'nullcompare': function(first, second) {
      return first == null ? second : first;
    },

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

  var $promise = {
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

    'translate': function(prom) {
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

  var moduleInits = [];

  var $module = function(name, creator) {
  	var module = {};

    var parts = name.split('.');
    var current = $g;
    for (var i = 0; i < parts.length - 1; ++i) {
      if (!current[parts[i]]) {
        current[parts[i]] = {};
      }
      current = current[parts[i]]
    }

    current[parts[parts.length - 1]] = module;

  	module.$init = function(cpromise) {
  	  moduleInits.push(cpromise);
  	};

    module.$newtypebuilder = function(kind) {
      return function(name, hasGenerics, alias, creator) {
        var buildType = function(n, args) {
          var tpe = new Function("return function " + n + "() {};")();
          creator.apply(tpe, args || []);

          if (kind == 'type') {
            tpe.prototype.toJSON = function() {
              return $t.nominalunwrap(this);
            };
          } else if (kind == 'struct') {
            // Stringify.
            tpe.prototype.Stringify = function(T) {
              var $this = this;
              return function() {
                // Special case JSON, as it uses an internal method.
                if (T == $a['$json']) {
                  return $promise.resolve(JSON.stringify($this.data));
                }

                return T.Get().then(function(resolved) {
                  return resolved.Stringify($t.any)($this.data);
                });
              };
            };

            // Parse.
            tpe.Parse = function(T) {
              return function(value) {
                // TODO: Validate the struct.

                // Special case JSON for performance, as it uses an internal method.
                if (T == $a['$json']) {
                  var created = new tpe();
                  created.data = JSON.parse($t.nominalunwrap(value));
                  return $promise.resolve(created);
                }

                return T.Get().then(function(resolved) {
                  return (resolved.Parse($t.any)(value)).then(function(parsed) {
                    var created = new tpe();
                    // TODO: *efficiently* unwrap internal nominal types.
                    created.data = JSON.parse(JSON.stringify($t.nominalunwrap(parsed)));
                    return $promise.resolve(created);
                  });
                });
              };
            };
          }

          return tpe;
        };

        if (hasGenerics) {
          module[name] = function(__genericargs) {
            var fullName = name;
            for (var i = 0; i < arguments.length; ++i) {
              fullName = fullName + '_' + arguments[i].name;
            }

            return buildType(fullName, arguments);
          };
        } else {
          module[name] = buildType(name);
        }

        if (alias) {
          $a[alias] = module[name];
        }
      };
    };

    module.$struct = module.$newtypebuilder('struct');
  	module.$class = module.$newtypebuilder('class');
  	module.$interface = module.$newtypebuilder('interface');
    module.$type = module.$newtypebuilder('type');

  	creator.call(module)
  };

  {{ range $idx, $kv := .Iter }}
  	{{ $kv.Value }}
  {{ end }}

  $g.executeWorkerMethod = function(token) {
    $global.onmessage = function(e) {
      if (!e.isTrusted) { return; }

      var data = e.data;
      if (data['token'] != token) {
        throw Error('Invalid token')
      }

      switch (data['action']) {
        case 'invoke':
          var methodId = data['method'];
          var arguments = data['arguments'];
          var method = $w[methodId];

          method.apply(null, arguments).then(function(result) {
            $global.postMessage({
              'result': result,
              'token': token
            });
            close();
          }).catch(function(reject) {
            $global.postMessage({
              'reject': reject,
              'token': token
            });
            close();
          });
          break
      }
    };
  };

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
          global.executeWorkerMethod(pair[1]);
        });
        return;
      }
    }

  };
  runWorker();
}
`
