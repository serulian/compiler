"use strict";
this.Serulian = function ($global) {
  var BOXED_DATA_PROPERTY = '$wrapped';
  var $__currentScriptSrc = null;
  if (typeof $global.document === 'object') {
    $__currentScriptSrc = $global.document.currentScript.src;
  }
  $global.__serulian_internal = {
    autoUnbox: function (k, v) {
      return $t.unbox(v);
    },
    autoBox: function (k, v) {
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
    },
  };
  var $g = {
  };
  var $a = {
  };
  var $w = {
  };
  var $it = function (name, typeIndex) {
    var tpe = new Function(("return function " + name) + "() {};")();
    tpe.$typeId = typeIndex;
    tpe.$typeref = function () {
      return {
        i: typeIndex,
      };
    };
    return tpe;
  };
  var $t = {
    any: $it('Any', 'any'),
    struct: $it('Struct', 'struct'),
    void: $it('Void', 'void'),
    null: $it('Null', 'null'),
    toESType: function (obj) {
      return {
      }.toString.call(obj).match(/\s([a-zA-Z]+)/)[1].toLowerCase();
    },
    ensureerror: function (rejected) {
      if (rejected instanceof Error) {
        return $a['wrappederror'].For(rejected);
      }
      return rejected;
    },
    markpromising: function (func) {
      func.$promising = true;
      return func;
    },
    functionName: function (func) {
      if (func.name) {
        return func.name;
      }
      var ret = func.toString();
      ret = ret.substr('function '.length);
      ret = ret.substr(0, ret.indexOf('('));
      return ret;
    },
    typeid: function (type) {
      return type.$typeId || $t.functionName(type);
    },
    buildDataForValue: function (value) {
      if (value == null) {
        return {
          v: null,
        };
      }
      if (value.constructor.$typeref) {
        return {
          v: $t.unbox(value),
          t: value.constructor.$typeref(),
        };
      } else {
        return {
          v: value,
        };
      }
    },
    buildValueFromData: function (data) {
      if (!data['t']) {
        return data['v'];
      }
      return $t.box(data['v'], $t.typeforref(data['t']));
    },
    unbox: function (instance) {
      if ((instance != null) && instance.hasOwnProperty(BOXED_DATA_PROPERTY)) {
        return instance[BOXED_DATA_PROPERTY];
      }
      return instance;
    },
    box: function (instance, type) {
      if (instance == null) {
        return null;
      }
      if (instance.constructor == type) {
        return instance;
      }
      return type.$box($t.unbox(instance));
    },
    fastbox: function (instance, type) {
      return type.$box(instance);
    },
    roottype: function (type) {
      if (type.$roottype) {
        return type.$roottype();
      }
      return type;
    },
    istype: function (value, type) {
      if ((type == $t.any) || ((value != null) && ((value.constructor == type) || (value instanceof type)))) {
        return true;
      }
      if (type == $t.struct) {
        var roottype = $t.roottype(value.constructor);
        return ((((roottype.$typekind == 'struct') || (roottype == Number)) || (roottype == String)) || (roottype == Boolean)) || (roottype == Object);
      }
      if ((type.$generic == $a['function']) && (typeof value == 'function')) {
        return value;
      }
      var targetKind = type.$typekind;
      switch (targetKind) {
        case 'struct':

        case 'type':

        case 'class':
          return false;

        case 'interface':
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
    cast: function (value, type, opt_allownull) {
      if (((value == null) && !opt_allownull) && (type != $t.any)) {
        throw Error('Cannot cast null value to ' + type.toString());
      }
      if ($t.istype(value, type)) {
        return value;
      }
      var targetKind = type.$typekind;
      switch (targetKind) {
        case 'struct':
          if (value.constructor == Object) {
            break;
          }
          throw Error((('Cannot cast ' + value.constructor.toString()) + ' to ') + type.toString());

        case 'class':

        case 'interface':
          throw Error((('Cannot cast ' + value.constructor.toString()) + ' to ') + type.toString());

        case 'type':
          if ($t.roottype(value.constructor) != $t.roottype(type)) {
            throw Error((('Cannot auto-box ' + value.constructor.toString()) + ' to ') + type.toString());
          }
          break;

        case undefined:
          throw Error((('Cannot cast ' + value.constructor.toString()) + ' to ') + type.toString());
      }
      if (type.$box) {
        return $t.box(value, type);
      }
      return value;
    },
    equals: function (left, right, type) {
      if (left === right) {
        return true;
      }
      if ((left == null) || (right == null)) {
        return false;
      }
      if (type.$equals) {
        return type.$equals($t.box(left, type), $t.box(right, type)).$wrapped;
      }
      return false;
    },
    ensurevalue: function (value, type, canBeNull, name) {
      if (value == null) {
        if (!canBeNull) {
          throw Error('Missing value for non-nullable field ' + name);
        }
        return;
      }
      var check = function (serutype, estype) {
        if ((type == $a[serutype]) || (type.$generic == $a[serutype])) {
          if ($t.toESType(value) != estype) {
            throw Error((((('Expected ' + serutype) + ' for field ') + name) + ', found: ') + $t.toESType(value));
          }
          return true;
        }
        return false;
      };
      if (check('string', 'string')) {
        return;
      }
      if (check('float64', 'number')) {
        return;
      }
      if (check('int', 'number')) {
        return;
      }
      if (check('bool', 'boolean')) {
        return;
      }
      if (check('slice', 'array')) {
        return;
      }
      if ($t.toESType(value) != 'object') {
        throw Error((('Expected object for field ' + name) + ', found: ') + $t.toESType(value));
      }
    },
    nativenew: function (type) {
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
    typeforref: function (typeref) {
      if (typeref['i']) {
        return $t[typeref['i']];
      }
      var parts = typeref['t'].split('.');
      var current = $g;
      for (var i = 0; i < parts.length; ++i) {
        current = current[parts[i]];
      }
      if (!typeref['g'].length) {
        return current;
      }
      var generics = typeref['g'].map(function (generic) {
        return $t.typeforref(generic);
      });
      return current.apply(current, generics);
    },
    uuid: function () {
      var buf = new Uint16Array(8);
      crypto.getRandomValues(buf);
      var S4 = function (num) {
        var ret = num.toString(16);
        while (ret.length < 4) {
          ret = "0" + ret;
        }
        return ret;
      };
      return ((((((((((S4(buf[0]) + S4(buf[1])) + "-") + S4(buf[2])) + "-") + S4(buf[3])) + "-") + S4(buf[4])) + "-") + S4(buf[5])) + S4(buf[6])) + S4(buf[7]);
    },
    defineStructField: function (structType, name, serializableName, typeref, opt_nominalRootType, opt_nullAllowed) {
      var field = {
        name: name,
        serializableName: serializableName,
        typeref: typeref,
        nominalRootTyperef: opt_nominalRootType || typeref,
        nullAllowed: opt_nullAllowed,
      };
      structType.$fields.push(field);
      Object.defineProperty(structType.prototype, name, {
        get: function () {
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
          return boxedData[name];
        },
        set: function (value) {
          this[BOXED_DATA_PROPERTY][name] = value;
        },
      });
    },
    workerwrap: function (methodId, f) {
      $w[methodId] = f;
      if (!$__currentScriptSrc) {
        return function () {
          var $this = this;
          var args = new Array(arguments.length);
          for (var i = 0; i < args.length; ++i) {
            args[i] = arguments[i];
          }
          var promise = new Promise(function (resolve, reject) {
            $global.setTimeout(function () {
              $promise.maybe(f.apply($this, args)).then(function (value) {
                resolve(value);
              }).catch(function (value) {
                reject(value);
              });
            }, 0);
          });
          return promise;
        };
      }
      return function () {
        var token = $t.uuid();
        var args = Array.prototype.map.call(arguments, $t.buildDataForValue);
        var promise = new Promise(function (resolve, reject) {
          var worker = new Worker(($__currentScriptSrc + "?__serulian_async_token=") + token);
          worker.onmessage = function (e) {
            if (!e.isTrusted) {
              worker.terminate();
              return;
            }
            var data = e.data;
            if (data['token'] != token) {
              return;
            }
            var value = $t.buildValueFromData(data['value']);
            var kind = data['kind'];
            if (kind == 'resolve') {
              resolve(value);
            } else {
              reject(value);
            }
            worker.terminate();
          };
          worker.postMessage({
            action: 'invoke',
            arguments: args,
            method: methodId,
            token: token,
          });
        });
        return promise;
      };
    },
    property: function (getter) {
      getter.$property = true;
      return getter;
    },
    nullableinvoke: function (obj, name, promising, args) {
      var found = obj != null ? obj[name] : null;
      if (found == null) {
        return promising ? $promise.resolve(null) : null;
      }
      var r = found.apply(obj, args);
      if (promising) {
        return $promise.maybe(r);
      } else {
        return r;
      }
    },
    dynamicaccess: function (obj, name, promising) {
      if ((obj == null) || (obj[name] == null)) {
        return promising ? $promise.resolve(null) : null;
      }
      var value = obj[name];
      if (typeof value == 'function') {
        if (value.$property) {
          var result = value.apply(obj, arguments);
          return promising ? $promise.maybe(result) : result;
        } else {
          var result = function () {
            return value.apply(obj, arguments);
          };
          return promising ? $promise.resolve(result) : result;
        }
      }
      return promising ? $promise.resolve(value) : value;
    },
    assertnotnull: function (value) {
      if (value == null) {
        throw Error('Value should not be null');
      }
      return value;
    },
    syncnullcompare: function (value, otherwise) {
      return value == null ? otherwise() : value;
    },
    asyncnullcompare: function (value, otherwise) {
      return value == null ? otherwise : value;
    },
    resourcehandler: function () {
      return {
        resources: {
        },
        bind: function (func, isAsync) {
          if (isAsync) {
            return this.bindasync(func);
          } else {
            return this.bindsync(func);
          }
        },
        bindsync: function (func) {
          var r = this;
          var f = function () {
            r.popall();
            return func.apply(this, arguments);
          };
          return f;
        },
        bindasync: function (func) {
          var r = this;
          var f = function (value) {
            var that = this;
            return r.popall().then(function (_) {
              func.call(that, value);
            });
          };
          return f;
        },
        pushr: function (value, name) {
          this.resources[name] = value;
        },
        popr: function (__names) {
          var handlers = [];
          for (var i = 0; i < arguments.length; ++i) {
            var name = arguments[i];
            if (this.resources[name]) {
              handlers.push(this.resources[name].Release());
              delete this.resources[name];
            }
          }
          return $promise.maybeall(handlers);
        },
        popall: function () {
          var handlers = [];
          var names = Object.keys(this.resources);
          for (var i = 0; i < names.length; ++i) {
            handlers.push(this.resources[names[i]].Release());
          }
          return $promise.maybeall(handlers);
        },
      };
    },
  };
  var $generator = {
    directempty: function () {
      var stream = {
        Next: function () {
          return $a['tuple']($t.any, $a['bool']).Build(null, false);
        },
      };
      return stream;
    },
    empty: function () {
      return $generator.directempty();
    },
    new: function (f, isAsync) {
      if (isAsync) {
        var stream = {
          $is: null,
          Next: function () {
            return $promise.new(function (resolve, reject) {
              if (stream.$is != null) {
                $promise.maybe(stream.$is.Next()).then(function (tuple) {
                  if ($t.unbox(tuple.Second)) {
                    resolve(tuple);
                  } else {
                    stream.$is = null;
                    $promise.maybe(stream.Next()).then(resolve, reject);
                  }
                }).catch(function (rejected) {
                  reject(rejected);
                });
                return;
              }
              var $yield = function (value) {
                resolve($a['tuple']($t.any, $a['bool']).Build(value, $t.fastbox(true, $a['bool'])));
              };
              var $done = function () {
                resolve($a['tuple']($t.any, $a['bool']).Build(null, $t.fastbox(false, $a['bool'])));
              };
              var $yieldin = function (ins) {
                stream.$is = ins;
                $promise.maybe(stream.Next()).then(resolve, reject);
              };
              f($yield, $yieldin, reject, $done);
            });
          },
        };
        return stream;
      } else {
        var stream = {
          $is: null,
          Next: function () {
            if (stream.$is != null) {
              var tuple = stream.$is.Next();
              if ($t.unbox(tuple.Second)) {
                return tuple;
              } else {
                stream.$is = null;
              }
            }
            var yielded = null;
            var $yield = function (value) {
              yielded = $a['tuple']($t.any, $a['bool']).Build(value, $t.fastbox(true, $a['bool']));
            };
            var $done = function () {
              yielded = $a['tuple']($t.any, $a['bool']).Build(null, $t.fastbox(false, $a['bool']));
            };
            var $yieldin = function (ins) {
              stream.$is = ins;
            };
            var $reject = function (rejected) {
              throw rejected;
            };
            f($yield, $yieldin, $reject, $done);
            if (stream.$is) {
              return stream.Next();
            } else {
              return yielded;
            }
          },
        };
        return stream;
      }
    },
  };
  var $promise = {
    all: function (promises) {
      return Promise.all(promises);
    },
    maybeall: function (results) {
      return Promise.all(results.map($promise.maybe));
    },
    maybe: function (r) {
      if (r && r.then) {
        return r;
      } else {
        return Promise.resolve(r);
      }
    },
    new: function (f) {
      return new Promise(f);
    },
    empty: function () {
      return Promise.resolve(null);
    },
    resolve: function (value) {
      return Promise.resolve(value);
    },
    reject: function (value) {
      return Promise.reject(value);
    },
    wrap: function (func) {
      return Promise.resolve(func());
    },
    shortcircuit: function (left, right) {
      if (left != right) {
        return $promise.resolve(left);
      }
    },
    translate: function (prom) {
      if (!prom.Then) {
        return prom;
      }
      return {
        then: function () {
          return prom.Then.apply(prom, arguments);
        },
        catch: function () {
          return prom.Catch.apply(prom, arguments);
        },
      };
    },
  };
  var moduleInits = [];
  var $module = function (moduleName, creator) {
    var module = {
    };
    var parts = moduleName.split('.');
    var current = $g;
    for (var i = 0; i < (parts.length - 1); ++i) {
      if (!current[parts[i]]) {
        current[parts[i]] = {
        };
      }
      current = current[parts[i]];
    }
    current[parts[parts.length - 1]] = module;
    var $newtypebuilder = function (kind) {
      return function (typeId, name, hasGenerics, alias, creator) {
        var buildType = function (fullTypeId, fullName, args) {
          var args = args || [];
          var tpe = new Function(("return function " + fullName) + "() {};")();
          tpe.$typeref = function () {
            if (!hasGenerics) {
              return {
                t: (moduleName + '.') + name,
              };
            }
            var generics = [];
            for (var i = 0; i < args.length; ++i) {
              generics.push(args[i].$typeref());
            }
            return {
              t: (moduleName + '.') + name,
              g: generics,
            };
          };
          tpe.$typeId = fullTypeId;
          tpe.$typekind = kind;
          creator.apply(tpe, args);
          if (kind == 'struct') {
            tpe.$box = function (data) {
              var instance = new tpe();
              instance[BOXED_DATA_PROPERTY] = data;
              instance.$lazychecked = {
              };
              return instance;
            };
            tpe.prototype.$markruntimecreated = function () {
              Object.defineProperty(this[BOXED_DATA_PROPERTY], '$runtimecreated', {
                enumerable: false,
                configurable: true,
                value: true,
              });
            };
            tpe.prototype.String = function () {
              return $t.fastbox(JSON.stringify(this, $global.__serulian_internal.autoUnbox, ' '), $a['string']);
            };
            tpe.prototype.Clone = function () {
              var instance = new tpe();
              if (Object.assign) {
                instance[BOXED_DATA_PROPERTY] = Object.assign({
                }, this[BOXED_DATA_PROPERTY]);
              } else {
                instance[BOXED_DATA_PROPERTY] = {
                };
                for (var key in this[BOXED_DATA_PROPERTY]) {
                  if (this[BOXED_DATA_PROPERTY].hasOwnProperty(key)) {
                    instance[BOXED_DATA_PROPERTY][key] = this[BOXED_DATA_PROPERTY][key];
                  }
                }
              }
              if (this[BOXED_DATA_PROPERTY].$runtimecreated) {
                instance.$markruntimecreated();
              }
              return instance;
            };
            tpe.prototype.Stringify = function (T) {
              var $this = this;
              return function () {
                if (T == $a['json']) {
                  return $promise.resolve($t.fastbox(JSON.stringify($this, $global.__serulian_internal.autoUnbox), $a['string']));
                }
                var mapped = $this.Mapping();
                return $promise.maybe(T.Get()).then(function (resolved) {
                  return resolved.Stringify(mapped);
                });
              };
            };
            tpe.Parse = function (T) {
              return function (value) {
                if (T == $a['json']) {
                  var parsed = JSON.parse($t.unbox(value));
                  var boxed = $t.fastbox(parsed, tpe);
                  var initPromise = $promise.resolve(boxed);
                  if (tpe.$initDefaults) {
                    initPromise = $promise.maybe(tpe.$initDefaults(boxed, false));
                  }
                  return initPromise.then(function () {
                    boxed.Mapping();
                    return boxed;
                  });
                }
                return $promise.maybe(T.Get()).then(function (resolved) {
                  return $promise.maybe(resolved.Parse(value)).then(function (parsed) {
                    return $promise.resolve($t.box(parsed, tpe));
                  });
                });
              };
            };
            tpe.$equals = function (left, right) {
              if (left === right) {
                return $t.fastbox(true, $a['bool']);
              }
              for (var i = 0; i < tpe.$fields.length; ++i) {
                var field = tpe.$fields[i];
                if (!$t.equals(left[BOXED_DATA_PROPERTY][field.serializableName], right[BOXED_DATA_PROPERTY][field.serializableName], field.typeref())) {
                  return $t.fastbox(false, $a['bool']);
                }
              }
              return $t.fastbox(true, $a['bool']);
            };
            tpe.prototype.Mapping = function () {
              if (this.$serucreated) {
                return $t.fastbox(this[BOXED_DATA_PROPERTY], $a['mapping']($t.any));
              } else {
                var $this = this;
                var mapped = {
                };
                tpe.$fields.forEach(function (field) {
                  mapped[field.serializableName] = $this[field.name];
                });
                return $t.fastbox(mapped, $a['mapping']($t.any));
              }
            };
          }
          return tpe;
        };
        if (hasGenerics) {
          module[name] = function genericType () {
            var fullName = name;
            var fullId = typeId;
            var generics = new Array(arguments.length);
            for (var i = 0; i < generics.length; ++i) {
              fullName = (fullName + '_') + $t.functionName(arguments[i]);
              if (i == 0) {
                fullId = fullId + '<';
              } else {
                fullId = fullId + ',';
              }
              fullId = fullId + arguments[i].$typeId;
              generics[i] = arguments[i];
            }
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
        if (alias) {
          $a[alias] = module[name];
        }
      };
    };
    module.$init = function (callback, fieldId, dependencyIds) {
      moduleInits.push({
        callback: callback,
        id: fieldId,
        depends: dependencyIds,
      });
    };
    module.$struct = $newtypebuilder('struct');
    module.$class = $newtypebuilder('class');
    module.$agent = $newtypebuilder('agent');
    module.$interface = $newtypebuilder('interface');
    module.$type = $newtypebuilder('type');
    creator.call(module);
  };
  $module('____testlib.basictypes', function () {
    var $static = this;
    this.$class('58998129', 'Tuple', true, 'tuple', function (T, Q) {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        instance.First = /*#null#*/null;
        instance.Second = /*#null#*/null;
        return instance;
      };
      $static.Build = function (first, second) {
        var tuple;
        tuple = /*#Tuple<T, Q>.new()#*/$g.____testlib.basictypes.Tuple(/*#T, Q>.new()#*/T, /*#Q>.new()#*/Q).new();
/*#tuple.First = first#*/        tuple.First = /*#first#*/first;
/*#tuple.Second = second#*/        tuple.Second = /*#second#*/second;
        return tuple;
      };
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
        };
        computed[((("Build|1|fd8bc7c9<58998129<" + $t.typeid(T)) + ",") + $t.typeid(Q)) + ">>"] = true;
        return this.$cachedtypesig = computed;
      };
    });

    this.$class('fd8bc7c9', 'Function', true, 'function', function (T) {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        return instance;
      };
      this.$typesig = function () {
        return {
        };
      };
    });

    this.$class('08e38881', 'IntStream', false, '$intstream', function () {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        instance.start = /*#0#*/$t.fastbox(/*#0#*/0, /*#0#*/$g.____testlib.basictypes.Integer);
        instance.end = /*#-1#*/$t.fastbox(/*#-1#*/-/*#-1#*/1, /*#-1#*/$g.____testlib.basictypes.Integer);
        instance.current = /*#0#*/$t.fastbox(/*#0#*/0, /*#0#*/$g.____testlib.basictypes.Integer);
        return instance;
      };
      $static.OverRange = function (start, end) {
        var s;
        s = /*#IntStream.new()#*/$g.____testlib.basictypes.IntStream.new();
/*#s.start = start#*/        s.start = /*#start#*/start;
/*#s.end = end#*/        s.end = /*#end#*/end;
/*#s.current = start#*/        s.current = /*#start#*/start;
        return s;
      };
      $instance.Next = function () {
        var $this = this;
        var t;
        var $current = 0;
        syncloop: while (true) {
          switch ($current) {
            case 0:
              if (/*#this.current <= this.end {#*/$this.current.$wrapped <= /*#this.end {#*/$this.end.$wrapped) /*#this.end {#*/{
                $current = 1;
                continue syncloop;
              } else {
                $current = 2;
                continue syncloop;
              }
              break;

            case 1:
              t = /*#Tuple<int, bool>.Build(this.current, true)#*/$g.____testlib.basictypes.Tuple(/*#int, bool>.Build(this.current, true)#*/$g.____testlib.basictypes.Integer, /*#bool>.Build(this.current, true)#*/$g.____testlib.basictypes.Boolean).Build(/*#this.current, true)#*/$this.current, /*#true)#*/$t.fastbox(/*#true)#*/true, /*#true)#*/$g.____testlib.basictypes.Boolean));
/*#this.current = this.current + 1#*/              $this.current = /*#this.current + 1#*/$t.fastbox(/*#this.current + 1#*/$this.current.$wrapped + /*#this.current + 1#*/1, /*#this.current + 1#*/$g.____testlib.basictypes.Integer);
              return t;

            case 2:
              return $g.____testlib.basictypes.Tuple($g.____testlib.basictypes.Integer, $g.____testlib.basictypes.Boolean).Build($this.current, $t.fastbox(false, $g.____testlib.basictypes.Boolean));

            default:
              return;
          }
        }
      };
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "OverRange|1|fd8bc7c9<08e38881>": true,
          "Next|2|fd8bc7c9<58998129<bb8d3aad,9706e8ab>>": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$class('db8daea9', 'Float64', false, 'float64', function () {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        return instance;
      };
      this.$typesig = function () {
        return {
        };
      };
    });

    this.$class('6af5f55b', 'List', true, 'list', function (T) {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        instance.items = /*#Array.new()#*/$t.nativenew(/*#Array.new()#*/$global.Array)();
        return instance;
      };
      $static.forArray = function (arr) {
        var l;
        l = /*#List<T>.new()#*/$g.____testlib.basictypes.List(/*#T>.new()#*/T).new();
/*#l.items = arr#*/        l.items = /*#arr#*/arr;
        return l;
      };
      $instance.Count = $t.property(function () {
        var $this = this;
        return $t.fastbox($this.items.length, $g.____testlib.basictypes.Integer);
      });
      $instance.$index = function (index) {
        var $this = this;
        return $t.cast($this.items[index.$wrapped], T, false);
      };
      $instance.$slice = function (startindex, endindex) {
        var $this = this;
        var end;
        var start;
        var $current = 0;
        syncloop: while (true) {
          switch ($current) {
            case 0:
              start = /*#startindex ?? 0#*/$t.syncnullcompare(/*#startindex ?? 0#*/startindex, function () {
                return $t.fastbox(0, $g.____testlib.basictypes.Integer);
              });
              end = /*#endindex ?? this.Count#*/$t.syncnullcompare(/*#endindex ?? this.Count#*/endindex, function () {
                return $this.Count();
              });
              if (/*#start < 0 {#*/start.$wrapped < /*#start < 0 {#*/0) /*#start < 0 {#*/{
                $current = 1;
                continue syncloop;
              } else {
                $current = 2;
                continue syncloop;
              }
              break;

            case 1:
/*#start = start + this.Count#*/              start = /*#start + this.Count#*/$t.fastbox(/*#start + this.Count#*/start.$wrapped + /*#this.Count#*/$this.Count().$wrapped, /*#start + this.Count#*/$g.____testlib.basictypes.Integer);
              $current = 2;
              continue syncloop;

            case 2:
              if (/*#end < 0 {#*/end.$wrapped < /*#end < 0 {#*/0) /*#end < 0 {#*/{
                $current = 3;
                continue syncloop;
              } else {
                $current = 4;
                continue syncloop;
              }
              break;

            case 3:
/*#end = end + this.Count#*/              end = /*#end + this.Count#*/$t.fastbox(/*#end + this.Count#*/end.$wrapped + /*#this.Count#*/$this.Count().$wrapped, /*#end + this.Count#*/$g.____testlib.basictypes.Integer);
              $current = 4;
              continue syncloop;

            case 4:
              if (/*#start >= end {#*/start.$wrapped >= /*#end {#*/end.$wrapped) /*#end {#*/{
                $current = 5;
                continue syncloop;
              } else {
                $current = 6;
                continue syncloop;
              }
              break;

            case 5:
              return $g.____testlib.basictypes.Slice(T).Empty();

            case 6:
              return $g.____testlib.basictypes.Slice(T).overArray($this.items.slice(start.$wrapped, end.$wrapped));

            default:
              return;
          }
        }
      };
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "Count|3|bb8d3aad": true,
        };
        computed[("index|4|fd8bc7c9<" + $t.typeid(T)) + ">"] = true;
        computed[("slice|4|fd8bc7c9<9b1698d6<" + $t.typeid(T)) + ">>"] = true;
        return this.$cachedtypesig = computed;
      };
    });

    this.$class('2e031d94', 'Map', true, 'map', function (T, Q) {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        instance.internalObject = /*#Object.new()#*/$t.nativenew(/*#Object.new()#*/$global.Object)();
        return instance;
      };
      $static.Empty = function () {
        return $g.____testlib.basictypes.Map(T, Q).new();
      };
      $static.forArrays = function (keys, values) {
        var $temp0;
        var $temp1;
        var i;
        var len;
        var map;
        var tKey;
        var $current = 0;
        syncloop: while (true) {
          switch ($current) {
            case 0:
              map = /*#Map<T, Q>.new()#*/$g.____testlib.basictypes.Map(/*#T, Q>.new()#*/T, /*#Q>.new()#*/Q).new();
              len = /*#int(keys.length)#*/$t.fastbox(/*#keys.length)#*/keys.length, /*#int(keys.length)#*/$g.____testlib.basictypes.Integer);
              $current = 1;
              continue syncloop;

            case 1:
              $temp1 = /*#0..(len - 1) {#*/$g.____testlib.basictypes.Integer.$range(/*#0..(len - 1) {#*/$t.fastbox(/*#0..(len - 1) {#*/0, /*#0..(len - 1) {#*/$g.____testlib.basictypes.Integer), /*#len - 1) {#*/$t.fastbox(/*#len - 1) {#*/len.$wrapped - /*#len - 1) {#*/1, /*#len - 1) {#*/$g.____testlib.basictypes.Integer));
              $current = 2;
              continue syncloop;

            case 2:
/*#i in 0..(len - 1) {#*/              $temp0 = /*#for i in 0..(len - 1) {#*/$temp1.Next();
/*#i in 0..(len - 1) {#*/              i = /*#i in 0..(len - 1) {#*/$temp0.First;
              if (/*#for i in 0..(len - 1) {#*/$temp0.Second.$wrapped) /*#for i in 0..(len - 1) {#*/{
                $current = 3;
                continue syncloop;
              } else {
                $current = 4;
                continue syncloop;
              }
              break;

            case 3:
              tKey = /*#keys[Number(i)].(T)#*/$t.cast(/*#keys[Number(i)].(T)#*/keys[/*#i)].(T)#*/i.$wrapped], /*#keys[Number(i)].(T)#*/T, /*#keys[Number(i)].(T)#*/false);
/*#map[tKey] = values[Number(i)].(Q)#*/              map.$setindex(/*#tKey] = values[Number(i)].(Q)#*/tKey, /*#values[Number(i)].(Q)#*/$t.cast(/*#values[Number(i)].(Q)#*/values[/*#i)].(Q)#*/i.$wrapped], /*#values[Number(i)].(Q)#*/Q, /*#values[Number(i)].(Q)#*/false));
              $current = 2;
              continue syncloop;

            case 4:
              return map;

            default:
              return;
          }
        }
      };
      $instance.Mapping = function () {
        var $this = this;
        return $t.fastbox($this.internalObject, $g.____testlib.basictypes.Mapping(Q));
      };
      $instance.$index = function (key) {
        var $this = this;
        var keyString;
        var value;
        var $current = 0;
        syncloop: while (true) {
          switch ($current) {
            case 0:
              keyString = /*#key.MapKey.String()#*/key.MapKey().String();
              value = /*#this.internalObject[NativeString(keyString)]#*/$this.internalObject[/*#keyString)]#*/keyString.$wrapped];
              if (/*#value is null {#*/value == /*#null {#*/null) /*#null {#*/{
                $current = 1;
                continue syncloop;
              } else {
                $current = 2;
                continue syncloop;
              }
              break;

            case 1:
              return null;

            case 2:
              return $t.cast(value, Q, false);

            default:
              return;
          }
        }
      };
      $instance.$setindex = function (key, value) {
        var $this = this;
        var keyString;
        keyString = /*#key.MapKey.String()#*/key.MapKey().String();
/*#this.internalObject[NativeString(keyString)] = value#*/        $this.internalObject[/*#keyString)] = value#*/keyString.$wrapped] = /*#value#*/value;
        return;
      };
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "setindex|4|fd8bc7c9<void>": true,
        };
        computed[((("Empty|1|fd8bc7c9<2e031d94<" + $t.typeid(T)) + ",") + $t.typeid(Q)) + ">>"] = true;
        computed[("Mapping|2|fd8bc7c9<ad6de9ce<" + $t.typeid(Q)) + ">>"] = true;
        computed[("index|4|fd8bc7c9<" + $t.typeid(Q)) + ">"] = true;
        return this.$cachedtypesig = computed;
      };
    });

    this.$class('e84fb0ce', 'JSON', false, 'json', function () {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        return instance;
      };
      $static.Get = function () {
        return $g.____testlib.basictypes.JSON.new();
      };
      $instance.Stringify = function (value) {
        var $this = this;
        return $t.fastbox($global.JSON.stringify(value.$wrapped, $t.dynamicaccess($global.__serulian_internal, 'autoUnbox', false)), $g.____testlib.basictypes.String);
      };
      $instance.Parse = function (value) {
        var $this = this;
        return $t.fastbox($global.JSON.parse(value.$wrapped, $t.dynamicaccess($global.__serulian_internal, 'autoBox', false)), $g.____testlib.basictypes.Mapping($t.any));
      };
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "Get|1|fd8bc7c9<e84fb0ce>": true,
          "Stringify|2|fd8bc7c9<268aa058>": true,
          "Parse|2|fd8bc7c9<ad6de9ce<any>>": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$interface('8ca0891a', 'Stringable', false, 'stringable', function () {
      var $static = this;
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "String|2|fd8bc7c9<268aa058>": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$interface('1d733667', 'Stream', true, 'stream', function (T) {
      var $static = this;
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
        };
        computed[("Next|2|fd8bc7c9<58998129<" + $t.typeid(T)) + ",9706e8ab>>"] = true;
        return this.$cachedtypesig = computed;
      };
    });

    this.$interface('0d4528b3', 'Streamable', true, 'streamable', function (T) {
      var $static = this;
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
        };
        computed[("Stream|2|fd8bc7c9<1d733667<" + $t.typeid(T)) + ">>"] = true;
        return this.$cachedtypesig = computed;
      };
    });

    this.$interface('b6ef544c', 'Error', false, 'error', function () {
      var $static = this;
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "Message|3|268aa058": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$interface('69a8a979', 'Awaitable', true, 'awaitable', function (T) {
      var $static = this;
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
        };
        computed[("Then|2|fd8bc7c9<69a8a979<" + $t.typeid(T)) + ">>"] = true;
        computed[("Catch|2|fd8bc7c9<69a8a979<" + $t.typeid(T)) + ">>"] = true;
        return this.$cachedtypesig = computed;
      };
    });

    this.$interface('8875fbad', 'Releasable', false, 'releasable', function () {
      var $static = this;
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "Release|2|fd8bc7c9<void>": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$interface('f06c2356', 'Mappable', false, 'mappable', function () {
      var $static = this;
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "MapKey|3|8ca0891a": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$interface('e4152f8e', 'Stringifier', false, '$stringifier', function () {
      var $static = this;
      $static.Get = function () {
        return $g.____testlib.basictypes.JSON.new();
      };
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "Get|1|fd8bc7c9<e4152f8e>": true,
          "Stringify|2|fd8bc7c9<268aa058>": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$interface('f5873a5f', 'Parser', false, '$parser', function () {
      var $static = this;
      $static.Get = function () {
        return $g.____testlib.basictypes.JSON.new();
      };
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "Get|1|fd8bc7c9<f5873a5f>": true,
          "Parse|2|fd8bc7c9<ad6de9ce<any>>": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$type('ad6de9ce', 'Mapping', true, 'mapping', function (T) {
      var $instance = this.prototype;
      var $static = this;
      this.$box = function ($wrapped) {
        var instance = new this();
        instance[BOXED_DATA_PROPERTY] = $wrapped;
        return instance;
      };
      this.$roottype = function () {
        return $global.Object;
      };
      $static.Empty = function () {
        return $t.fastbox($t.nativenew($global.Object)(), $g.____testlib.basictypes.Mapping(T));
      };
      $static.overObject = function (obj) {
        return $t.fastbox(obj, $g.____testlib.basictypes.Mapping(T));
      };
      $instance.Keys = $t.property(function () {
        var $this = this;
        return $g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).overArray($global.Object.keys($this.$wrapped));
      });
      $instance.$index = function (key) {
        var $this = this;
        var value;
        var $current = 0;
        syncloop: while (true) {
          switch ($current) {
            case 0:
              value = /*#this)[NativeString(key)]#*/$this.$wrapped[/*#key)]#*/key.$wrapped];
              if (/*#value is null {#*/value == /*#null {#*/null) /*#null {#*/{
                $current = 1;
                continue syncloop;
              } else {
                $current = 2;
                continue syncloop;
              }
              break;

            case 1:
              return null;

            case 2:
              return $t.cast(value, T, false);

            default:
              return;
          }
        }
      };
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "Keys|3|9b1698d6<268aa058>": true,
        };
        computed[("Empty|1|fd8bc7c9<ad6de9ce<" + $t.typeid(T)) + ">>"] = true;
        computed[("index|4|fd8bc7c9<" + $t.typeid(T)) + ">"] = true;
        return this.$cachedtypesig = computed;
      };
    });

    this.$type('9b1698d6', 'Slice', true, 'slice', function (T) {
      var $instance = this.prototype;
      var $static = this;
      this.$box = function ($wrapped) {
        var instance = new this();
        instance[BOXED_DATA_PROPERTY] = $wrapped;
        return instance;
      };
      this.$roottype = function () {
        return $global.Array;
      };
      $static.Empty = function () {
        return $t.fastbox($t.nativenew($global.Array)(), $g.____testlib.basictypes.Slice(T));
      };
      $static.overArray = function (arr) {
        return $t.fastbox(arr, $g.____testlib.basictypes.Slice(T));
      };
      $instance.$index = function (index) {
        var $this = this;
        return $t.cast($this.$wrapped[index.$wrapped], T, false);
      };
      $instance.$slice = function (startindex, endindex) {
        var $this = this;
        var end;
        var start;
        var $current = 0;
        syncloop: while (true) {
          switch ($current) {
            case 0:
              start = /*#startindex ?? 0#*/$t.syncnullcompare(/*#startindex ?? 0#*/startindex, function () {
                return $t.fastbox(0, $g.____testlib.basictypes.Integer);
              });
              end = /*#endindex ?? this.Length#*/$t.syncnullcompare(/*#endindex ?? this.Length#*/endindex, function () {
                return $this.Length();
              });
              if (/*#start < 0 {#*/start.$wrapped < /*#start < 0 {#*/0) /*#start < 0 {#*/{
                $current = 1;
                continue syncloop;
              } else {
                $current = 2;
                continue syncloop;
              }
              break;

            case 1:
/*#start = start + this.Length#*/              start = /*#start + this.Length#*/$t.fastbox(/*#start + this.Length#*/start.$wrapped + /*#this.Length#*/$this.Length().$wrapped, /*#start + this.Length#*/$g.____testlib.basictypes.Integer);
              $current = 2;
              continue syncloop;

            case 2:
              if (/*#end < 0 {#*/end.$wrapped < /*#end < 0 {#*/0) /*#end < 0 {#*/{
                $current = 3;
                continue syncloop;
              } else {
                $current = 4;
                continue syncloop;
              }
              break;

            case 3:
/*#end = end + this.Length#*/              end = /*#end + this.Length#*/$t.fastbox(/*#end + this.Length#*/end.$wrapped + /*#this.Length#*/$this.Length().$wrapped, /*#end + this.Length#*/$g.____testlib.basictypes.Integer);
              $current = 4;
              continue syncloop;

            case 4:
              if (/*#start >= end {#*/start.$wrapped >= /*#end {#*/end.$wrapped) /*#end {#*/{
                $current = 5;
                continue syncloop;
              } else {
                $current = 6;
                continue syncloop;
              }
              break;

            case 5:
              return $g.____testlib.basictypes.Slice(T).Empty();

            case 6:
              return $g.____testlib.basictypes.Slice(T).overArray($this.$wrapped.slice(start.$wrapped, end.$wrapped));

            default:
              return;
          }
        }
      };
      $instance.Length = $t.property(function () {
        var $this = this;
        return $t.fastbox($this.$wrapped.length, $g.____testlib.basictypes.Integer);
      });
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "Length|3|bb8d3aad": true,
        };
        computed[("Empty|1|fd8bc7c9<9b1698d6<" + $t.typeid(T)) + ">>"] = true;
        computed[("index|4|fd8bc7c9<" + $t.typeid(T)) + ">"] = true;
        computed[("slice|4|fd8bc7c9<9b1698d6<" + $t.typeid(T)) + ">>"] = true;
        return this.$cachedtypesig = computed;
      };
    });

    this.$type('bb8d3aad', 'Integer', false, 'int', function () {
      var $instance = this.prototype;
      var $static = this;
      this.$box = function ($wrapped) {
        var instance = new this();
        instance[BOXED_DATA_PROPERTY] = $wrapped;
        return instance;
      };
      this.$roottype = function () {
        return $global.Number;
      };
      $static.$range = function (start, end) {
        return $g.____testlib.basictypes.IntStream.OverRange(start, end);
      };
      $static.$compare = function (left, right) {
        return $t.fastbox(left.$wrapped - right.$wrapped, $g.____testlib.basictypes.Integer);
      };
      $static.$equals = function (left, right) {
        return $t.box(left.$wrapped == right.$wrapped, $g.____testlib.basictypes.Boolean);
      };
      $static.$plus = function (left, right) {
        return $t.fastbox(left.$wrapped + right.$wrapped, $g.____testlib.basictypes.Integer);
      };
      $static.$minus = function (left, right) {
        return $t.fastbox(left.$wrapped - right.$wrapped, $g.____testlib.basictypes.Integer);
      };
      $instance.Release = function () {
        var $this = this;
        return;
      };
      $instance.MapKey = $t.property(function () {
        var $this = this;
        return $this;
      });
      $instance.String = function () {
        var $this = this;
        return $t.fastbox($this.$wrapped.toString(), $g.____testlib.basictypes.String);
      };
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "range|4|fd8bc7c9<1d733667<bb8d3aad>>": true,
          "compare|4|fd8bc7c9<bb8d3aad>": true,
          "equals|4|fd8bc7c9<9706e8ab>": true,
          "plus|4|fd8bc7c9<bb8d3aad>": true,
          "minus|4|fd8bc7c9<bb8d3aad>": true,
          "Release|2|fd8bc7c9<void>": true,
          "MapKey|3|8ca0891a": true,
          "String|2|fd8bc7c9<268aa058>": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$type('9706e8ab', 'Boolean', false, 'bool', function () {
      var $instance = this.prototype;
      var $static = this;
      this.$box = function ($wrapped) {
        var instance = new this();
        instance[BOXED_DATA_PROPERTY] = $wrapped;
        return instance;
      };
      this.$roottype = function () {
        return $global.Boolean;
      };
      $static.$compare = function (left, right) {
        var $current = 0;
        syncloop: while (true) {
          switch ($current) {
            case 0:
              if (/*#left == right {#*/$g.____testlib.basictypes.Boolean.$equals(/*#left == right {#*/left, /*#right {#*/right).$wrapped) /*#right {#*/{
                $current = 1;
                continue syncloop;
              } else {
                $current = 2;
                continue syncloop;
              }
              break;

            case 1:
              return $t.fastbox(0, $g.____testlib.basictypes.Integer);

            case 2:
              return $t.fastbox(-1, $g.____testlib.basictypes.Integer);

            default:
              return;
          }
        }
      };
      $static.$equals = function (left, right) {
        return $t.box(left.$wrapped == right.$wrapped, $g.____testlib.basictypes.Boolean);
      };
      $instance.String = function () {
        var $this = this;
        return $t.fastbox($this.$wrapped.toString(), $g.____testlib.basictypes.String);
      };
      $instance.MapKey = $t.property(function () {
        var $this = this;
        return $this;
      });
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "compare|4|fd8bc7c9<bb8d3aad>": true,
          "equals|4|fd8bc7c9<9706e8ab>": true,
          "String|2|fd8bc7c9<268aa058>": true,
          "MapKey|3|8ca0891a": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$type('268aa058', 'String', false, 'string', function () {
      var $instance = this.prototype;
      var $static = this;
      this.$box = function ($wrapped) {
        var instance = new this();
        instance[BOXED_DATA_PROPERTY] = $wrapped;
        return instance;
      };
      this.$roottype = function () {
        return $global.String;
      };
      $instance.String = function () {
        var $this = this;
        return $this;
      };
      $static.$equals = function (first, second) {
        return $t.box(first.$wrapped == second.$wrapped, $g.____testlib.basictypes.Boolean);
      };
      $static.$plus = function (first, second) {
        return $t.fastbox(first.$wrapped + second.$wrapped, $g.____testlib.basictypes.String);
      };
      $instance.MapKey = $t.property(function () {
        var $this = this;
        return $this;
      });
      $instance.Length = $t.property(function () {
        var $this = this;
        return $t.fastbox($this.$wrapped.length, $g.____testlib.basictypes.Integer);
      });
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "String|2|fd8bc7c9<268aa058>": true,
          "equals|4|fd8bc7c9<9706e8ab>": true,
          "plus|4|fd8bc7c9<268aa058>": true,
          "MapKey|3|8ca0891a": true,
          "Length|3|bb8d3aad": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    this.$type('443094ae', 'WrappedError', false, 'wrappederror', function () {
      var $instance = this.prototype;
      var $static = this;
      this.$box = function ($wrapped) {
        var instance = new this();
        instance[BOXED_DATA_PROPERTY] = $wrapped;
        return instance;
      };
      this.$roottype = function () {
        return $global.Error;
      };
      $static.For = function (err) {
        return $t.fastbox(err, $g.____testlib.basictypes.WrappedError);
      };
      $instance.Message = $t.property(function () {
        var $this = this;
        return $t.fastbox($this.$wrapped.message, $g.____testlib.basictypes.String);
      });
      this.$typesig = function () {
        if (this.$cachedtypesig) {
          return this.$cachedtypesig;
        }
        var computed = {
          "For|1|fd8bc7c9<443094ae>": true,
          "Message|3|268aa058": true,
        };
        return this.$cachedtypesig = computed;
      };
    });

    $static.formatTemplateString = function (pieces, values) {
      var $temp0;
      var $temp1;
      var i;
      var result;
      var $current = 0;
      syncloop: while (true) {
        switch ($current) {
          case 0:
            result = /*#''#*/$t.fastbox(/*#''#*/'', /*#''#*/$g.____testlib.basictypes.String);
            $current = 1;
            continue syncloop;

          case 1:
            $temp1 = /*#0 .. (pieces.Length - 1) {#*/$g.____testlib.basictypes.Integer.$range(/*#0 .. (pieces.Length - 1) {#*/$t.fastbox(/*#0 .. (pieces.Length - 1) {#*/0, /*#0 .. (pieces.Length - 1) {#*/$g.____testlib.basictypes.Integer), /*#pieces.Length - 1) {#*/$t.fastbox(/*#pieces.Length - 1) {#*/pieces.Length().$wrapped - /*#pieces.Length - 1) {#*/1, /*#pieces.Length - 1) {#*/$g.____testlib.basictypes.Integer));
            $current = 2;
            continue syncloop;

          case 2:
/*#i in 0 .. (pieces.Length - 1) {#*/            $temp0 = /*#for i in 0 .. (pieces.Length - 1) {#*/$temp1.Next();
/*#i in 0 .. (pieces.Length - 1) {#*/            i = /*#i in 0 .. (pieces.Length - 1) {#*/$temp0.First;
            if (/*#for i in 0 .. (pieces.Length - 1) {#*/$temp0.Second.$wrapped) /*#for i in 0 .. (pieces.Length - 1) {#*/{
              $current = 3;
              continue syncloop;
            } else {
              $current = 6;
              continue syncloop;
            }
            break;

          case 3:
/*#result = result + pieces[i]#*/            result = /*#result + pieces[i]#*/$g.____testlib.basictypes.String.$plus(/*#result + pieces[i]#*/result, /*#pieces[i]#*/pieces.$index(/*#i]#*/i));
            if (/*#i < values.Length {#*/i.$wrapped < /*#values.Length {#*/values.Length().$wrapped) /*#values.Length {#*/{
              $current = 4;
              continue syncloop;
            } else {
              $current = 5;
              continue syncloop;
            }
            break;

          case 4:
/*#result = result + values[i].String()#*/            result = /*#result + values[i].String()#*/$g.____testlib.basictypes.String.$plus(/*#result + values[i].String()#*/result, /*#values[i].String()#*/values.$index(/*#i].String()#*/i).String());
            $current = 5;
            continue syncloop;

          case 5:
            $current = 2;
            continue syncloop;

          case 6:
            return result;

          default:
            return;
        }
      }
    };
    $static.MapStream = function (T, Q) {
      var $f = function (stream, mapper) {
        var $result;
        var $temp0;
        var $temp1;
        var item;
        var $current = 0;
        var $continue = function ($yield, $yieldin, $reject, $done) {
          while (true) {
            switch ($current) {
              case 0:
                $current = 1;
                $continue($yield, $yieldin, $reject, $done);
                return;

              case 1:
                $temp1 = /*#stream {#*/stream;
                $current = 2;
                $continue($yield, $yieldin, $reject, $done);
                return;

              case 2:
/*#item in stream {#*/                $promise.maybe(/*#for item in stream {#*/$temp1.Next()).then(/*#for item in stream {#*/function (/*#for item in stream {#*/$result0) /*#for item in stream {#*/{
                  $temp0 = /*#item in stream {#*/$result0;
                  $result = /*#item in stream {#*/$temp0;
                  $current = 3;
                  $continue($yield, $yieldin, $reject, $done);
                  return;
                }).catch(function (err) {
                  throw err;
                });
                return;

              case 3:
/*#item in stream {#*/                item = /*#item in stream {#*/$temp0.First;
                if (/*#for item in stream {#*/$temp0.Second.$wrapped) /*#for item in stream {#*/{
                  $current = 4;
                  $continue($yield, $yieldin, $reject, $done);
                  return;
                } else {
                  $current = 6;
                  $continue($yield, $yieldin, $reject, $done);
                  return;
                }
                break;

              case 4:
                $yield(/*#mapper(item)#*/mapper(/*#item)#*/item));
                $current = 5;
                return;

              case 5:
                $current = 2;
                $continue($yield, $yieldin, $reject, $done);
                return;

              default:
                $done();
                return;
            }
          }
        };
        return $generator.new($continue, true);
      };
      return $f;
    };
  });
  $module('basic', function () {
    var $static = this;
    $static.DoSomething = function () {
      var bar;
      var foo;
      foo = /*#1#*/$t.fastbox(/*#1#*/1, /*#1#*/$g.____testlib.basictypes.Integer);
      bar = /*#'hi there!'#*/$t.fastbox(/*#'hi there!'#*/'hi there!', /*#'hi there!'#*/$g.____testlib.basictypes.String);
      return $t.fastbox((foo.$wrapped == 2) && $g.____testlib.basictypes.String.$equals(bar, $t.fastbox('hello world', $g.____testlib.basictypes.String)).$wrapped, $g.____testlib.basictypes.Boolean);
    };
  });
  $g.$executeWorkerMethod = function (token) {
    $global.onmessage = function (e) {
      if (!e.isTrusted) {
        $global.close();
        return;
      }
      var data = e.data;
      if (data['token'] != token) {
        throw Error('Invalid token');
      }
      switch (data['action']) {
        case 'invoke':
          var methodId = data['method'];
          var method = $w[methodId];
          var args = data['arguments'].map($t.buildValueFromData);
          var send = function (kind) {
            return function (value) {
              var message = {
                token: token,
                value: $t.buildDataForValue(value),
                kind: kind,
              };
              try {
                $global.postMessage(message);
              } catch (e) {
                if (kind == 'reject') {
                  throw value;
                } else {
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
  var buildPromises = function (items) {
    var seen = {
    };
    var result = [];
    var itemsById = {
    };
    items.forEach(function (item) {
      itemsById[item.id] = item;
    });
    items.forEach(function visit (item) {
      if (seen[item.id]) {
        return;
      }
      seen[item.id] = true;
      item.depends.forEach(function (depId) {
        visit(itemsById[depId]);
      });
      item['promise'] = item['callback']();
    });
    return items.map(function (item) {
      if (!item.depends.length) {
        return item['promise'];
      }
      var current = $promise.resolve();
      item.depends.forEach(function (depId) {
        current = current.then(function (resolved) {
          return itemsById[depId]['promise'];
        });
      });
      return current.then(function (resolved) {
        return item['promise'];
      });
    });
  };
  return $promise.all(buildPromises(moduleInits)).then(function () {
    return $g;
  });
}(this);
if (typeof importScripts === 'function') {
  var runWorker = function () {
    var search = location.search;
    if (!search || (search[0] != '?')) {
      return;
    }
    var searchPairs = search.substr(1).split('&');
    if (searchPairs.length < 1) {
      return;
    }
    for (var i = 0; i < searchPairs.length; ++i) {
      var pair = searchPairs[i].split('=');
      if (pair[0] == '__serulian_async_token') {
        this.Serulian.then(function (global) {
          global.$executeWorkerMethod(pair[1]);
        });
        return;
      }
    }
    close();
  };
  runWorker();
}

