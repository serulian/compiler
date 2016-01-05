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
    'nativenew': function(type, name, nominalType, promisenoop, promisewrap) {
      return function () {
        var newInstance = Object.create(type.prototype);
        newInstance = type.apply(newInstance, arguments) || newInstance;
        return promisewrap ? $promise.wrap(function () {
          return newInstance;
        }) : newInstance;
      };
    },

    'extension': function(obj, name, nominalType, promisenoop, promisewrap) {
       if (obj == null) {
          return $t.dynamicaccess(obj, name, null, promisenoop, promisewrap);
       }

       var func = function() {
          var args = [];
          args.push(obj);
          args.push.apply(args, arguments);
          return nominalType[name].apply(null, args)
       };

       return promisewrap ? $promise.wrap(func) : func;
    },

    'dynamicaccess': function(obj, name, nominalType, promisenoop, promisewrap) {
      if (obj == null || obj[name] == null) {
        return promisenoop ? $promise.wrap(function() { return null; }) : null;
      }

      var value = obj[name];
      if (typeof value == 'function') {
        if (promisewrap) {
          return $promise.wrap(function() {
            return value.apply(obj, arguments);
          });
        } else {
          return function() {
            return value.apply(obj, arguments);
          };
        }
      }

      return value
    },

  	'sm': function(caller) {
  		return {
        resources: {},
  			current: 0,
  			returnValue: undefined,
  			error: undefined,
  			next: caller,
        pushr: function(name, value) {
          this.resources[name] = value;
        },
        popr: function(name) {
          if (this.resources[name]) {
            this.resources[name].Release()
            delete this.resources[name];
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
  			var continueFunc = function() {
  				if (statemachine.returnValue !== undefined) {
            statemachine.popall();
  					resolve(statemachine.returnValue);
  					return;
  				}

  				if (statemachine.error !== undefined) {
            statemachine.popall();
  					reject(statemachine.error);
  					return;					
  				}

  				if (statemachine.current < 0) {
  					return;
  				}

  				statemachine.next(continueFunc);				
			};
			continueFunc();
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

  	'wrap': function(func) {
  		return Promise.resolve(func());
  	}
  };

  var moduleInits = [];

  var $module = function(name, creator) {
  	var module = {};
  	module.$init = function(cpromise) {
  	  moduleInits.push(cpromise());
  	};

  	module.$class = function(name, creator) {
  		var cls = function() {};
  		creator.call(cls);
  		module[name] = cls;
  	};

  	module.$interface = function(name, creator) {
  		var cls = function() {};
  		creator.call(cls);
  		module[name] = cls;
  	};

    module.$type = function(name, creator) {
      var cls = function() {};
      creator.call(cls);
      module[name] = cls;
    };

  	creator.call(module)
  	$g[name] = module;
  };

  {{ range $idx, $kv := .Iter }}
  	{{ $kv.Value }}
  {{ end }}

  return $promise.all(moduleInits).then(function() {
  	return $g;
  });
})(window)
`
