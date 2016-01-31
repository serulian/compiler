// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

// Note: nativenew is based on http://www.bennadel.com/blog/2291-invoking-a-native-javascript-constructor-using-call-or-apply.htm

// runtimeTemplate contains all the necessary code for wrapping generated modules into a complete Serulian
// runtime bundle.
const runtimeTemplate = `
window.Serulian = (function($global) {
  var $g = {};
  var $t = {
    'cast': function(value, type) {
      // TODO: implement cast checking.
      return value
    },

    'nativenew': function(type) {
      return function () {
        var newInstance = Object.create(type.prototype);
        newInstance = type.apply(newInstance, arguments) || newInstance;
        return newInstance;
      };
    },

    'property': function(isExtension, getter, opt_setter) {
      var count = isExtension ? 2 : 1;
      var f = function() {
        if (arguments.length == count) {
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
      return function(name, hasGenerics, creator) {
        if (hasGenerics) {
          module[name] = function(genericargs) {
            var tpe = function() {};
            creator.apply(tpe, arguments);
            return tpe;
          };
        } else {
          var tpe = new Function("return function " + name + "() {};")();
          creator.call(tpe);
          module[name] = tpe;
        }
      };
    };

  	module.$class = module.$newtypebuilder('class');
  	module.$interface = module.$newtypebuilder('interface');

    module.$type = function(name, creator) {
      var cls = function() {};
      creator.call(cls);
      module[name] = cls;
    };

  	creator.call(module)
  };

  {{ range $idx, $kv := .Iter }}
  	{{ $kv.Value }}
  {{ end }}

  return $promise.all(moduleInits).then(function() {
  	return $g;
  });
})(window)
`
