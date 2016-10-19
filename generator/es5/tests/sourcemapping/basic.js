"use strict";
this.Serulian = function ($global) {
  var BOXED_DATA_PROPERTY = '$wrapped';
  var $__currentScriptSrc = null;
  if (typeof $global.document === 'object') {
    $__currentScriptSrc = $global.document.currentScript.src;
  }
  $global.__serulian_internal = {
    autoUnbox: function (k, v) {
      if (v != null) {
        return $t.unbox(v);
      }
      return v;
    },
    autoBox: function (k, v) {
      if (v == null) {
        return v;
      }
      var typeName = $t.toESType(v);
      switch (typeName) {
        case 'object':
          if (k != '') {
            return $t.box(v, $a.mapping($t.any));
          }
          break;

        case 'array':
          return $t.box(v, $a.slice($t.any));

        case 'boolean':
          return $t.box(v, $a.bool);

        case 'string':
          return $t.box(v, $a.string);

        case 'number':
          if (Math.ceil(v) == v) {
            return $t.box(v, $a.int);
          }
          return $t.box(v, $a.float64);
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
    box: function (instance, type, opt_external) {
      if (instance == null) {
        return null;
      }
      return type.$box($t.unbox(instance), opt_external);
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
      if ((value == null) && !opt_allownull) {
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
      }
      if (type.$box) {
        return $t.box(value, type);
      }
      return value;
    },
    createtypesig: function () {
      var sig = {
      };
      for (var i = 0; i < arguments.length; ++i) {
        var entry = arguments[i];
        var key = (((entry[0] + ':') + entry[1]) + ':') + JSON.stringify(entry[2]);
        sig[key] = true;
      }
      return sig;
    },
    equals: function (left, right, type) {
      if (left === right) {
        return $promise.resolve($t.box(true, $a['bool']));
      }
      if ((left == null) || (right == null)) {
        return $promise.resolve($t.box(false, $a['bool']));
      }
      if ($t.toESType(left) != 'object') {
        return $promise.resolve($t.box(left === right, $a['bool']));
      }
      if (type.$equals) {
        return type.$equals($t.box(left, type), $t.box(right, type));
      }
      return $promise.resolve($t.box(false, $a['bool']));
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
    defineStructField: function (structType, name, serializableName, typeref, isBoxed, opt_nominalRootType, opt_nullAllowed) {
      structType.$fields.push({
        name: name,
        serializableName: serializableName,
        typeref: typeref,
        nominalRootTyperef: opt_nominalRootType || typeref,
        isBoxed: isBoxed,
        nullAllowed: opt_nullAllowed,
      });
      Object.defineProperty(structType.prototype, name, {
        get: function () {
          return this[BOXED_DATA_PROPERTY][name];
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
          var args = arguments;
          var promise = new Promise(function (resolve, reject) {
            $global.setTimeout(function () {
              f.apply($this, args).then(function (value) {
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
      return found.apply(obj, args);
    },
    dynamicaccess: function (obj, name) {
      if ((obj == null) || (obj[name] == null)) {
        return $promise.resolve(null);
      }
      var value = obj[name];
      if (typeof value == 'function') {
        if (value.$property) {
          return value.apply(obj, arguments);
        } else {
          return $promise.resolve(function () {
            return value.apply(obj, arguments);
          });
        }
      }
      return $promise.resolve(value);
    },
    assertnotnull: function (value) {
      if (value == null) {
        throw Error('Value should not be null');
      }
      return value;
    },
    nullcompare: function (value, otherwise) {
      return value == null ? otherwise : value;
    },
    resourcehandler: function () {
      return {
        resources: {
        },
        bind: function (func) {
          if (func.$__resourcebound) {
            return func;
          }
          var r = this;
          var f = function () {
            r.popall();
            return func.apply(this, arguments);
          };
          f.$__resourcebound = true;
          return f;
        },
        pushr: function (value, name) {
          this.resources[name] = value;
        },
        popr: function (names) {
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
        popall: function () {
          for (var name in this.resources) {
            if (this.resources.hasOwnProperty(name)) {
              this.resources[name].Release();
            }
          }
        },
      };
    },
  };
  var $generator = {
    directempty: function () {
      var stream = {
        Next: function () {
          return $promise.new(function (resolve, reject) {
            $a['tuple']($t.any, $a['bool']).Build(null, false).then(resolve);
          });
        },
      };
      return stream;
    },
    empty: function () {
      return $promise.resolve($generator.directempty());
    },
    new: function (f) {
      var stream = {
        $is: null,
        Next: function () {
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
              $a['tuple']($t.any, $a['bool']).Build(value, true).then(resolve);
            };
            var $done = function () {
              $a['tuple']($t.any, $a['bool']).Build(null, false).then(resolve);
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
    },
  };
  var $promise = {
    all: function (promises) {
      return Promise.all(promises);
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
      return function (name, hasGenerics, alias, creator) {
        var buildType = function (n, args) {
          var args = args || [];
          var tpe = new Function(("return function " + n) + "() {};")();
          tpe.$typeref = function () {
            var generics = [];
            for (var i = 0; i < args.length; ++i) {
              generics.push(args[i].$typeref());
            }
            return {
              t: (moduleName + '.') + name,
              g: generics,
            };
          };
          tpe.$typekind = kind;
          creator.apply(tpe, args);
          if (kind == 'struct') {
            tpe.$box = function (data) {
              var instance = new tpe();
              instance[BOXED_DATA_PROPERTY] = data;
              instance.$lazychecked = {
              };
              instance.$unboxed = true;
              tpe.$fields.forEach(function (field) {
                Object.defineProperty(instance, field.name, {
                  get: function () {
                    if (!this.$lazychecked[field.name]) {
                      $t.ensurevalue($t.unbox(this[BOXED_DATA_PROPERTY][field.serializableName]), field.nominalRootTyperef(), field.nullAllowed, field.name);
                      this.$lazychecked[field.name] = true;
                    }
                    if (field.isBoxed) {
                      return $t.box(this[BOXED_DATA_PROPERTY][field.serializableName], field.typeref());
                    } else {
                      return this[BOXED_DATA_PROPERTY][field.serializableName];
                    }
                  },
                });
              });
              return instance;
            };
            tpe.prototype.String = function () {
              return $promise.resolve($t.box(JSON.stringify(this, $global.__serulian_internal.autoUnbox, ' '), $a['string']));
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
              return $promise.resolve(instance);
            };
            tpe.prototype.Stringify = function (T) {
              var $this = this;
              return function () {
                if (T == $a['json']) {
                  return $promise.resolve($t.box(JSON.stringify($this, $global.__serulian_internal.autoUnbox), $a['string']));
                }
                return $this.Mapping().then(function (mapped) {
                  return T.Get().then(function (resolved) {
                    return resolved.Stringify(mapped);
                  });
                });
              };
            };
            tpe.Parse = function (T) {
              return function (value) {
                if (T == $a['json']) {
                  var parsed = JSON.parse($t.unbox(value));
                  var boxed = $t.box(parsed, tpe, true);
                  return boxed.Mapping().then(function () {
                    return $promise.resolve(boxed);
                  });
                }
                return T.Get().then(function (resolved) {
                  return resolved.Parse(value).then(function (parsed) {
                    return $promise.resolve($t.box(parsed, tpe, true));
                  });
                });
              };
            };
            tpe.$equals = function (left, right) {
              if (left === right) {
                return $promise.resolve($t.box(true, $a['bool']));
              }
              var promises = [];
              tpe.$fields.forEach(function (field) {
                promises.push($t.equals(left[BOXED_DATA_PROPERTY][field.serializableName], right[BOXED_DATA_PROPERTY][field.serializableName], field.typeref()));
              });
              return Promise.all(promises).then(function (values) {
                for (var i = 0; i < values.length; i++) {
                  if (!$t.unbox(values[i])) {
                    return values[i];
                  }
                }
                return $t.box(true, $a['bool']);
              });
            };
            tpe.prototype.Mapping = function () {
              if (this.$unboxed) {
                var $this = this;
                var mapped = {
                };
                tpe.$fields.forEach(function (field) {
                  mapped[field.serializableName] = $this[field.name];
                });
                return $promise.resolve($t.box(mapped, $a['mapping']($t.any)));
              } else {
                return $promise.resolve($t.box(this[BOXED_DATA_PROPERTY], $a['mapping']($t.any)));
              }
            };
          }
          return tpe;
        };
        if (hasGenerics) {
          module[name] = function genericType (__genericargs) {
            var fullName = name;
            for (var i = 0; i < arguments.length; ++i) {
              fullName = (fullName + '_') + arguments[i].name;
            }
            var tpe = buildType(name, arguments);
            tpe.$generic = genericType;
            return tpe;
          };
        } else {
          module[name] = buildType(name);
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
    module.$interface = $newtypebuilder('interface');
    module.$type = $newtypebuilder('type');
    creator.call(module);
  };
  $module('____testlib.basictypes', function () {
    var $static = this;
    this.$class('Tuple', true, 'tuple', function (T, Q) {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        instance.First = /*#null#*/null;
        instance.Second = /*#null#*/null;
        return $promise.resolve(instance);
      };
      $static.Build = function (first, second) {
        var $result;
        var tuple;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#Tuple<T, Q>.new()#*/                $g.____testlib.basictypes.Tuple(/*#T, Q>.new()#*/T, /*#Q>.new()#*/Q).new().then(/*#Q>.new()#*/function (/*#Q>.new()#*/$result0) /*#Q>.new()#*/{
                  $result = /*#.new()#*/$result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                tuple = $result;
/*#tuple.First = first#*/                tuple.First = /*#first#*/first;
/*#tuple.Second = second#*/                tuple.Second = /*#second#*/second;
                $resolve(/*#tuple#*/tuple);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['Build', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Tuple(T, Q)).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Tuple(T, Q)).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$class('Function', true, 'function', function (T) {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        return $promise.resolve(instance);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Function(T)).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$class('IntStream', false, '', function () {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        instance.start = /*#0#*/$t.box(/*#0#*/0, /*#0#*/$g.____testlib.basictypes.Integer);
        instance.end = /*#-1#*/$t.box(/*#-1#*/-/*#-1#*/1, /*#-1#*/$g.____testlib.basictypes.Integer);
        instance.current = /*#0#*/$t.box(/*#0#*/0, /*#0#*/$g.____testlib.basictypes.Integer);
        return $promise.resolve(instance);
      };
      $static.OverRange = function (start, end) {
        var $result;
        var s;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#IntStream.new()#*/                $g.____testlib.basictypes.IntStream.new().then(/*#IntStream.new()#*/function (/*#IntStream.new()#*/$result0) /*#IntStream.new()#*/{
                  $result = /*#.new()#*/$result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                s = $result;
/*#s.start = start#*/                s.start = /*#start#*/start;
/*#s.end = end#*/                s.end = /*#end#*/end;
/*#s.current = start#*/                s.current = /*#start#*/start;
                $resolve(/*#s#*/s);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      $instance.Next = function () {
        var $this = this;
        var $result;
        var t;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#this.current <= this.end {#*/                $g.____testlib.basictypes.Integer.$compare(/*#this.current <= this.end {#*/$this.current, /*#this.end {#*/$this.end).then(/*#this.end {#*/function (/*#this.end {#*/$result0) /*#this.end {#*/{
                  $result = /*#this.current <= this.end {#*/$t.unbox(/*#this.current <= this.end {#*/$result0) <= /*#this.current <= this.end {#*/0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                if ($result) {
                  $current = 2;
                  continue;
                } else {
                  $current = 5;
                  continue;
                }
                break;

              case 2:
/*#Tuple<int, bool>.Build(this.current, true)#*/                $g.____testlib.basictypes.Tuple(/*#int, bool>.Build(this.current, true)#*/$g.____testlib.basictypes.Integer, /*#bool>.Build(this.current, true)#*/$g.____testlib.basictypes.Boolean).Build(/*#this.current, true)#*/$this.current, /*#true)#*/$t.box(/*#true)#*/true, /*#true)#*/$g.____testlib.basictypes.Boolean)).then(/*#true)#*/function (/*#true)#*/$result0) /*#true)#*/{
                  $result = /*#.Build(this.current, true)#*/$result0;
                  $current = 3;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 3:
                t = $result;
/*#this.current + 1#*/                $g.____testlib.basictypes.Integer.$plus(/*#this.current + 1#*/$this.current, /*#1#*/$t.box(/*#1#*/1, /*#1#*/$g.____testlib.basictypes.Integer)).then(/*#1#*/function (/*#1#*/$result0) /*#1#*/{
                  $result = /*#this.current = this.current + 1#*/$this.current = /*#this.current + 1#*/$result0;
                  $current = 4;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 4:
                $resolve(/*#t#*/t);
                return;

              case 5:
/*#Tuple<int, bool>.Build(this.current, false)#*/                $g.____testlib.basictypes.Tuple(/*#int, bool>.Build(this.current, false)#*/$g.____testlib.basictypes.Integer, /*#bool>.Build(this.current, false)#*/$g.____testlib.basictypes.Boolean).Build(/*#this.current, false)#*/$this.current, /*#false)#*/$t.box(/*#false)#*/false, /*#false)#*/$g.____testlib.basictypes.Boolean)).then(/*#false)#*/function (/*#false)#*/$result0) /*#false)#*/{
                  $result = /*#.Build(this.current, false)#*/$result0;
                  $current = 6;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 6:
                $resolve($result);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['OverRange', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.IntStream).$typeref()], ['Next', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Tuple($g.____testlib.basictypes.Integer, $g.____testlib.basictypes.Boolean)).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.IntStream).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$class('Float64', false, 'float64', function () {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        return $promise.resolve(instance);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Float64).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$class('List', true, 'list', function (T) {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        instance.items = /*#Array.new()#*/$t.nativenew(/*#Array.new()#*/$global.Array)();
        return $promise.resolve(instance);
      };
      $static.forArray = function (arr) {
        var $result;
        var l;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#List<T>.new()#*/                $g.____testlib.basictypes.List(/*#T>.new()#*/T).new().then(/*#T>.new()#*/function (/*#T>.new()#*/$result0) /*#T>.new()#*/{
                  $result = /*#.new()#*/$result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                l = $result;
/*#l.items = arr#*/                l.items = /*#arr#*/arr;
                $resolve(/*#l#*/l);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      $instance.Count = $t.property(function () {
        var $this = this;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#int(this.items.length)#*/$t.box(/*#this.items.length)#*/$this.items.length, /*#int(this.items.length)#*/$g.____testlib.basictypes.Integer));
          return;
        };
        return $promise.new($continue);
      });
      $instance.$index = function (index) {
        var $this = this;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#[&index].(T)#*/$t.cast(/*#this.items[&index].(T)#*/$this.items[/*#index].(T)#*/$t.unbox(/*#index].(T)#*/index)], /*#[&index].(T)#*/T, /*#[&index].(T)#*/false));
          return;
        };
        return $promise.new($continue);
      };
      $instance.$slice = function (startindex, endindex) {
        var $this = this;
        var $result;
        var end;
        var start;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#startindex ?? 0#*/                $promise.resolve(/*#startindex ?? 0#*/startindex).then(/*#startindex ?? 0#*/function (/*#startindex ?? 0#*/$result0) /*#startindex ?? 0#*/{
                  $result = /*#startindex ?? 0#*/$t.nullcompare(/*#startindex ?? 0#*/$result0, /*#0#*/$t.box(/*#0#*/0, /*#0#*/$g.____testlib.basictypes.Integer));
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                start = $result;
/*#endindex ?? this.Count#*/                $promise.resolve(/*#endindex ?? this.Count#*/endindex).then(/*#endindex ?? this.Count#*/function (/*#endindex ?? this.Count#*/$result0) /*#endindex ?? this.Count#*/{
                  return /*#this.Count#*/(/*#this.Count#*/$promise.shortcircuit(/*#this.Count#*/$result0, /*#endindex ?? this.Count#*/null) || /*#this.Count#*/$this.Count()).then(/*#this.Count#*/function (/*#this.Count#*/$result1) /*#this.Count#*/{
                    $result = /*#endindex ?? this.Count#*/$t.nullcompare(/*#endindex ?? this.Count#*/$result0, /*#this.Count#*/$result1);
                    $current = 2;
                    $continue($resolve, $reject);
                    return;
                  });
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 2:
                end = $result;
/*#start < 0 {#*/                $g.____testlib.basictypes.Integer.$compare(/*#start < 0 {#*/start, /*#0 {#*/$t.box(/*#0 {#*/0, /*#0 {#*/$g.____testlib.basictypes.Integer)).then(/*#0 {#*/function (/*#0 {#*/$result0) /*#0 {#*/{
                  $result = /*#start < 0 {#*/$t.unbox(/*#start < 0 {#*/$result0) < /*#start < 0 {#*/0;
                  $current = 3;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 3:
                if ($result) {
                  $current = 4;
                  continue;
                } else {
                  $current = 6;
                  continue;
                }
                break;

              case 4:
/*#this.Count#*/                $this.Count().then(/*#this.Count#*/function (/*#this.Count#*/$result1) /*#this.Count#*/{
                  return /*#start + this.Count#*/$g.____testlib.basictypes.Integer.$plus(/*#start + this.Count#*/start, /*#this.Count#*/$result1).then(/*#this.Count#*/function (/*#this.Count#*/$result0) /*#this.Count#*/{
                    start = /*#start + this.Count#*/$result0;
                    $result = /*#start = start + this.Count#*/start;
                    $current = 5;
                    $continue($resolve, $reject);
                    return;
                  });
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 5:
                $current = 6;
                continue;

              case 6:
/*#end < 0 {#*/                $g.____testlib.basictypes.Integer.$compare(/*#end < 0 {#*/end, /*#0 {#*/$t.box(/*#0 {#*/0, /*#0 {#*/$g.____testlib.basictypes.Integer)).then(/*#0 {#*/function (/*#0 {#*/$result0) /*#0 {#*/{
                  $result = /*#end < 0 {#*/$t.unbox(/*#end < 0 {#*/$result0) < /*#end < 0 {#*/0;
                  $current = 7;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 7:
                if ($result) {
                  $current = 8;
                  continue;
                } else {
                  $current = 10;
                  continue;
                }
                break;

              case 8:
/*#this.Count#*/                $this.Count().then(/*#this.Count#*/function (/*#this.Count#*/$result1) /*#this.Count#*/{
                  return /*#end + this.Count#*/$g.____testlib.basictypes.Integer.$plus(/*#end + this.Count#*/end, /*#this.Count#*/$result1).then(/*#this.Count#*/function (/*#this.Count#*/$result0) /*#this.Count#*/{
                    end = /*#end + this.Count#*/$result0;
                    $result = /*#end = end + this.Count#*/end;
                    $current = 9;
                    $continue($resolve, $reject);
                    return;
                  });
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 9:
                $current = 10;
                continue;

              case 10:
/*#start >= end {#*/                $g.____testlib.basictypes.Integer.$compare(/*#start >= end {#*/start, /*#end {#*/end).then(/*#end {#*/function (/*#end {#*/$result0) /*#end {#*/{
                  $result = /*#start >= end {#*/$t.unbox(/*#start >= end {#*/$result0) >= /*#start >= end {#*/0;
                  $current = 11;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 11:
                if ($result) {
                  $current = 12;
                  continue;
                } else {
                  $current = 14;
                  continue;
                }
                break;

              case 12:
/*#Slice<T>.Empty()#*/                $g.____testlib.basictypes.Slice(/*#T>.Empty()#*/T).Empty().then(/*#T>.Empty()#*/function (/*#T>.Empty()#*/$result0) /*#T>.Empty()#*/{
                  $result = /*#.Empty()#*/$result0;
                  $current = 13;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 13:
                $resolve($result);
                return;

              case 14:
/*#Slice<T>.overArray(this.items.slice(Number(start), Number(end)))#*/                $g.____testlib.basictypes.Slice(/*#T>.overArray(this.items.slice(Number(start), Number(end)))#*/T).overArray(/*#this.items.slice(Number(start), Number(end)))#*/$this.items.slice(/*#Number(start), Number(end)))#*/$t.unbox(/*#start), Number(end)))#*/start), /*#Number(end)))#*/$t.unbox(/*#end)))#*/end))).then(/*#end)))#*/function (/*#end)))#*/$result0) /*#end)))#*/{
                  $result = /*#.overArray(this.items.slice(Number(start), Number(end)))#*/$result0;
                  $current = 15;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 15:
                $resolve($result);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['forArray', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.List(T)).$typeref()], ['Count', 3, $g.____testlib.basictypes.Integer.$typeref()], ['index', 4, $g.____testlib.basictypes.Function(T).$typeref()], ['slice', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Slice(T)).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.List(T)).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$class('Map', true, 'map', function (T, Q) {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        instance.internalObject = /*#Object.new()#*/$t.nativenew(/*#Object.new()#*/$global.Object)();
        return $promise.resolve(instance);
      };
      $static.forArrays = function (keys, values) {
        var $result;
        var $temp0;
        var $temp1;
        var i;
        var len;
        var map;
        var tKey;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#Map<T, Q>.new()#*/                $g.____testlib.basictypes.Map(/*#T, Q>.new()#*/T, /*#Q>.new()#*/Q).new().then(/*#Q>.new()#*/function (/*#Q>.new()#*/$result0) /*#Q>.new()#*/{
                  $result = /*#.new()#*/$result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                map = $result;
                len = /*#int(keys.length)#*/$t.box(/*#keys.length)#*/keys.length, /*#int(keys.length)#*/$g.____testlib.basictypes.Integer);
                $current = 2;
                continue;

              case 2:
/*#len - 1) {#*/                $g.____testlib.basictypes.Integer.$minus(/*#len - 1) {#*/len, /*#1) {#*/$t.box(/*#1) {#*/1, /*#1) {#*/$g.____testlib.basictypes.Integer)).then(/*#1) {#*/function (/*#1) {#*/$result1) /*#1) {#*/{
                  return /*#0..(len - 1) {#*/$g.____testlib.basictypes.Integer.$range(/*#0..(len - 1) {#*/$t.box(/*#0..(len - 1) {#*/0, /*#0..(len - 1) {#*/$g.____testlib.basictypes.Integer), /*#len - 1) {#*/$result1).then(/*#len - 1) {#*/function (/*#len - 1) {#*/$result0) /*#len - 1) {#*/{
                    $result = /*#0..(len - 1) {#*/$result0;
                    $current = 3;
                    $continue($resolve, $reject);
                    return;
                  });
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 3:
                $temp1 = $result;
                $current = 4;
                continue;

              case 4:
/*#for i in 0..(len - 1) {#*/                $temp1.Next().then(/*#for i in 0..(len - 1) {#*/function (/*#for i in 0..(len - 1) {#*/$result0) /*#for i in 0..(len - 1) {#*/{
                  $temp0 = /*#i in 0..(len - 1) {#*/$result0;
                  $result = /*#i in 0..(len - 1) {#*/$temp0;
                  $current = 5;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 5:
/*#i in 0..(len - 1) {#*/                i = /*#i in 0..(len - 1) {#*/$temp0.First;
                if (/*#for i in 0..(len - 1) {#*/$t.unbox(/*#for i in 0..(len - 1) {#*/$temp0.Second)) /*#for i in 0..(len - 1) {#*/{
                  $current = 6;
                  continue;
                } else {
                  $current = 8;
                  continue;
                }
                break;

              case 6:
                tKey = /*#[Number(i)].(T)#*/$t.cast(/*#keys[Number(i)].(T)#*/keys[/*#Number(i)].(T)#*/$t.unbox(/*#i)].(T)#*/i)], /*#[Number(i)].(T)#*/T, /*#[Number(i)].(T)#*/false);
/*#map[tKey] = values[Number(i)].(Q)#*/                map.$setindex(/*#tKey] = values[Number(i)].(Q)#*/tKey, /*#[Number(i)].(Q)#*/$t.cast(/*#values[Number(i)].(Q)#*/values[/*#Number(i)].(Q)#*/$t.unbox(/*#i)].(Q)#*/i)], /*#[Number(i)].(Q)#*/Q, /*#[Number(i)].(Q)#*/false)).then(/*#[Number(i)].(Q)#*/function (/*#[Number(i)].(Q)#*/$result0) /*#[Number(i)].(Q)#*/{
                  $result = /*#map[tKey] = values[Number(i)].(Q)#*/$result0;
                  $current = 7;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 7:
                $current = 4;
                continue;

              case 8:
                $resolve(/*#map#*/map);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      $instance.$index = function (key) {
        var $this = this;
        var $result;
        var keyString;
        var value;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#key.MapKey.String()#*/                key.MapKey().then(/*#key.MapKey.String()#*/function (/*#key.MapKey.String()#*/$result1) /*#key.MapKey.String()#*/{
                  return /*#key.MapKey.String()#*/$result1.String().then(/*#key.MapKey.String()#*/function (/*#key.MapKey.String()#*/$result0) /*#key.MapKey.String()#*/{
                    $result = /*#.String()#*/$result0;
                    $current = 1;
                    $continue($resolve, $reject);
                    return;
                  });
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                keyString = $result;
                value = /*#this.internalObject[NativeString(keyString)]#*/$this.internalObject[/*#NativeString(keyString)]#*/$t.unbox(/*#keyString)]#*/keyString)];
                if (/*#value is null {#*/value == /*#null {#*/null) /*#null {#*/{
                  $current = 2;
                  continue;
                } else {
                  $current = 3;
                  continue;
                }
                break;

              case 2:
                $resolve(/*#null#*/null);
                return;

              case 3:
                $resolve(/*#value.(Q)#*/$t.cast(/*#value.(Q)#*/value, /*#value.(Q)#*/Q, /*#value.(Q)#*/false));
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      $instance.$setindex = function (key, value) {
        var $this = this;
        var $result;
        var keyString;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#key.MapKey.String()#*/                key.MapKey().then(/*#key.MapKey.String()#*/function (/*#key.MapKey.String()#*/$result1) /*#key.MapKey.String()#*/{
                  return /*#key.MapKey.String()#*/$result1.String().then(/*#key.MapKey.String()#*/function (/*#key.MapKey.String()#*/$result0) /*#key.MapKey.String()#*/{
                    $result = /*#.String()#*/$result0;
                    $current = 1;
                    $continue($resolve, $reject);
                    return;
                  });
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                keyString = $result;
/*#this.internalObject[NativeString(keyString)] = value#*/                $this.internalObject[/*#NativeString(keyString)] = value#*/$t.unbox(/*#keyString)] = value#*/keyString)] = /*#value#*/value;
                $resolve();
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['forArrays', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Map(T, Q)).$typeref()], ['index', 4, $g.____testlib.basictypes.Function(Q).$typeref()], ['setindex', 4, $g.____testlib.basictypes.Function($t.void).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Map(T, Q)).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$class('JSON', false, 'json', function () {
      var $static = this;
      var $instance = this.prototype;
      $static.new = function () {
        var instance = new $static();
        return $promise.resolve(instance);
      };
      $static.Get = function () {
        var $result;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#JSON.new()#*/                $g.____testlib.basictypes.JSON.new().then(/*#JSON.new()#*/function (/*#JSON.new()#*/$result0) /*#JSON.new()#*/{
                  $result = /*#.new()#*/$result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                $resolve($result);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      $instance.Stringify = function (value) {
        var $this = this;
        var $result;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#Internal.autoUnbox))#*/                $t.dynamicaccess(/*#Internal.autoUnbox))#*/$global.__serulian_internal, /*#Internal.autoUnbox))#*/'autoUnbox').then(/*#Internal.autoUnbox))#*/function (/*#Internal.autoUnbox))#*/$result0) /*#Internal.autoUnbox))#*/{
                  $result = /*#string(NativeJSON.stringify(Object(value), Internal.autoUnbox))#*/$t.box(/*#NativeJSON.stringify(Object(value), Internal.autoUnbox))#*/$global.JSON.stringify(/*#Object(value), Internal.autoUnbox))#*/$t.unbox(/*#value), Internal.autoUnbox))#*/value), /*#Internal.autoUnbox))#*/$result0), /*#string(NativeJSON.stringify(Object(value), Internal.autoUnbox))#*/$g.____testlib.basictypes.String);
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                $resolve($result);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      $instance.Parse = function (value) {
        var $this = this;
        var $result;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#Internal.autoBox))#*/                $t.dynamicaccess(/*#Internal.autoBox))#*/$global.__serulian_internal, /*#Internal.autoBox))#*/'autoBox').then(/*#Internal.autoBox))#*/function (/*#Internal.autoBox))#*/$result0) /*#Internal.autoBox))#*/{
                  $result = /*#<any>(NativeJSON.parse(NativeString(value), Internal.autoBox))#*/$t.box(/*#NativeJSON.parse(NativeString(value), Internal.autoBox))#*/$global.JSON.parse(/*#NativeString(value), Internal.autoBox))#*/$t.unbox(/*#value), Internal.autoBox))#*/value), /*#Internal.autoBox))#*/$result0), /*#<any>(NativeJSON.parse(NativeString(value), Internal.autoBox))#*/$g.____testlib.basictypes.Mapping(/*#<any>(NativeJSON.parse(NativeString(value), Internal.autoBox))#*/$t.any));
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                $resolve($result);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['Get', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.JSON).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Parse', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.JSON).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$interface('Stringable', false, 'stringable', function () {
      var $static = this;
      this.$typesig = function () {
        var computed = $t.createtypesig(['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$interface('Stream', true, 'stream', function (T) {
      var $static = this;
      this.$typesig = function () {
        var computed = $t.createtypesig(['Next', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Tuple(T, $g.____testlib.basictypes.Boolean)).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$interface('Streamable', true, 'streamable', function (T) {
      var $static = this;
      this.$typesig = function () {
        var computed = $t.createtypesig(['Stream', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Stream(T)).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$interface('Error', false, 'error', function () {
      var $static = this;
      this.$typesig = function () {
        var computed = $t.createtypesig(['Message', 3, $g.____testlib.basictypes.String.$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$interface('Awaitable', true, 'awaitable', function (T) {
      var $static = this;
      this.$typesig = function () {
        var computed = $t.createtypesig(['Then', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Awaitable(T)).$typeref()], ['Catch', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Awaitable(T)).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$interface('Releasable', false, 'releasable', function () {
      var $static = this;
      this.$typesig = function () {
        var computed = $t.createtypesig(['Release', 2, $g.____testlib.basictypes.Function($t.void).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$interface('Mappable', false, 'mappable', function () {
      var $static = this;
      this.$typesig = function () {
        var computed = $t.createtypesig(['MapKey', 3, $g.____testlib.basictypes.Stringable.$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$interface('Stringifier', false, '$stringifier', function () {
      var $static = this;
      $static.Get = function () {
        var $result;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#JSON.new()#*/                $g.____testlib.basictypes.JSON.new().then(/*#JSON.new()#*/function (/*#JSON.new()#*/$result0) /*#JSON.new()#*/{
                  $result = /*#.new()#*/$result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                $resolve($result);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['Get', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Stringifier).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$interface('Parser', false, '$parser', function () {
      var $static = this;
      $static.Get = function () {
        var $result;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#JSON.new()#*/                $g.____testlib.basictypes.JSON.new().then(/*#JSON.new()#*/function (/*#JSON.new()#*/$result0) /*#JSON.new()#*/{
                  $result = /*#.new()#*/$result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                $resolve($result);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['Get', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Parser).$typeref()], ['Parse', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$type('Mapping', true, 'mapping', function (T) {
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
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#<T>(Object.new())#*/$t.box(/*#Object.new())#*/$t.nativenew(/*#Object.new())#*/$global.Object)(), /*#<T>(Object.new())#*/$g.____testlib.basictypes.Mapping(/*#<T>(Object.new())#*/T)));
          return;
        };
        return $promise.new($continue);
      };
      $static.overObject = function (obj) {
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#<T>(obj)#*/$t.box(/*#obj)#*/obj, /*#<T>(obj)#*/$g.____testlib.basictypes.Mapping(/*#<T>(obj)#*/T)));
          return;
        };
        return $promise.new($continue);
      };
      $instance.Keys = $t.property(function () {
        var $this = this;
        var $result;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#Slice<string>.overArray(Object.keys(Object(this)))#*/                $g.____testlib.basictypes.Slice(/*#string>.overArray(Object.keys(Object(this)))#*/$g.____testlib.basictypes.String).overArray(/*#Object.keys(Object(this)))#*/$global.Object.keys(/*#Object(this)))#*/$t.unbox(/*#this)))#*/$this))).then(/*#this)))#*/function (/*#this)))#*/$result0) /*#this)))#*/{
                  $result = /*#.overArray(Object.keys(Object(this)))#*/$result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                $resolve($result);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      });
      $instance.$index = function (key) {
        var $this = this;
        var value;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
                value = /*#Object(this)[NativeString(key)]#*/$t.unbox(/*#this)[NativeString(key)]#*/$this)[/*#NativeString(key)]#*/$t.unbox(/*#key)]#*/key)];
                if (/*#value is null {#*/value == /*#null {#*/null) /*#null {#*/{
                  $current = 1;
                  continue;
                } else {
                  $current = 2;
                  continue;
                }
                break;

              case 1:
                $resolve(/*#null#*/null);
                return;

              case 2:
                $resolve(/*#value.(T)#*/$t.cast(/*#value.(T)#*/value, /*#value.(T)#*/T, /*#value.(T)#*/false));
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['Empty', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping(T)).$typeref()], ['overObject', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping(T)).$typeref()], ['Keys', 3, $g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).$typeref()], ['index', 4, $g.____testlib.basictypes.Function(T).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$type('Slice', true, 'slice', function (T) {
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
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#<T>(Array.new())#*/$t.box(/*#Array.new())#*/$t.nativenew(/*#Array.new())#*/$global.Array)(), /*#<T>(Array.new())#*/$g.____testlib.basictypes.Slice(/*#<T>(Array.new())#*/T)));
          return;
        };
        return $promise.new($continue);
      };
      $static.overArray = function (arr) {
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#<T>(arr)#*/$t.box(/*#arr)#*/arr, /*#<T>(arr)#*/$g.____testlib.basictypes.Slice(/*#<T>(arr)#*/T)));
          return;
        };
        return $promise.new($continue);
      };
      $instance.$index = function (index) {
        var $this = this;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#[Number(index)].(T)#*/$t.cast(/*#Array(this)[Number(index)].(T)#*/$t.unbox(/*#this)[Number(index)].(T)#*/$this)[/*#Number(index)].(T)#*/$t.unbox(/*#index)].(T)#*/index)], /*#[Number(index)].(T)#*/T, /*#[Number(index)].(T)#*/false));
          return;
        };
        return $promise.new($continue);
      };
      $instance.Length = $t.property(function () {
        var $this = this;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#int(Array(this).length)#*/$t.box(/*#Array(this).length)#*/$t.unbox(/*#this).length)#*/$this).length, /*#int(Array(this).length)#*/$g.____testlib.basictypes.Integer));
          return;
        };
        return $promise.new($continue);
      });
      this.$typesig = function () {
        var computed = $t.createtypesig(['Empty', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Slice(T)).$typeref()], ['overArray', 1, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Slice(T)).$typeref()], ['index', 4, $g.____testlib.basictypes.Function(T).$typeref()], ['Length', 3, $g.____testlib.basictypes.Integer.$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$type('Integer', false, 'int', function () {
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
        var $result;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#IntStream.OverRange(start, end)#*/                $g.____testlib.basictypes.IntStream.OverRange(/*#start, end)#*/start, /*#end)#*/end).then(/*#end)#*/function (/*#end)#*/$result0) /*#end)#*/{
                  $result = /*#.OverRange(start, end)#*/$result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                $resolve($result);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      $static.$compare = function (left, right) {
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#Integer(Number(left) - Number(right))#*/$t.box(/*#Number(left) - Number(right))#*/$t.unbox(/*#left) - Number(right))#*/left) - /*#Number(right))#*/$t.unbox(/*#right))#*/right), /*#Integer(Number(left) - Number(right))#*/$g.____testlib.basictypes.Integer));
          return;
        };
        return $promise.new($continue);
      };
      $static.$equals = function (left, right) {
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#Boolean(Number(left) == Number(right))#*/$t.box(/*#Number(left) == Number(right))#*/$t.unbox(/*#left) == Number(right))#*/left) == /*#Number(right))#*/$t.unbox(/*#right))#*/right), /*#Boolean(Number(left) == Number(right))#*/$g.____testlib.basictypes.Boolean));
          return;
        };
        return $promise.new($continue);
      };
      $static.$plus = function (left, right) {
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#Integer(Number(left) + Number(right))#*/$t.box(/*#Number(left) + Number(right))#*/$t.unbox(/*#left) + Number(right))#*/left) + /*#Number(right))#*/$t.unbox(/*#right))#*/right), /*#Integer(Number(left) + Number(right))#*/$g.____testlib.basictypes.Integer));
          return;
        };
        return $promise.new($continue);
      };
      $static.$minus = function (left, right) {
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#Integer(Number(left) - Number(right))#*/$t.box(/*#Number(left) - Number(right))#*/$t.unbox(/*#left) - Number(right))#*/left) - /*#Number(right))#*/$t.unbox(/*#right))#*/right), /*#Integer(Number(left) - Number(right))#*/$g.____testlib.basictypes.Integer));
          return;
        };
        return $promise.new($continue);
      };
      $instance.Release = function () {
        var $this = this;
        return $promise.empty();
      };
      $instance.MapKey = $t.property(function () {
        var $this = this;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#this }#*/$this);
          return;
        };
        return $promise.new($continue);
      });
      $instance.String = function () {
        var $this = this;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#String(Number(this).toString())#*/$t.box(/*#Number(this).toString())#*/$t.unbox(/*#this).toString())#*/$this).toString(), /*#String(Number(this).toString())#*/$g.____testlib.basictypes.String));
          return;
        };
        return $promise.new($continue);
      };
      this.$typesig = function () {
        var computed = $t.createtypesig(['range', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Stream($g.____testlib.basictypes.Integer)).$typeref()], ['compare', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Integer).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['plus', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Integer).$typeref()], ['minus', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Integer).$typeref()], ['Release', 2, $g.____testlib.basictypes.Function($t.void).$typeref()], ['MapKey', 3, $g.____testlib.basictypes.Stringable.$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$type('Boolean', false, 'bool', function () {
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
        var $result;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          while (true) {
            switch ($current) {
              case 0:
/*#left == right {#*/                $g.____testlib.basictypes.Boolean.$equals(/*#left == right {#*/left, /*#right {#*/right).then(/*#right {#*/function (/*#right {#*/$result0) /*#right {#*/{
                  $result = /*#if left == right {#*/$t.unbox(/*#left == right {#*/$result0);
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 1:
                if ($result) {
                  $current = 2;
                  continue;
                } else {
                  $current = 3;
                  continue;
                }
                break;

              case 2:
                $resolve(/*#0#*/$t.box(/*#0#*/0, /*#0#*/$g.____testlib.basictypes.Integer));
                return;

              case 3:
                $resolve(/*#-1#*/$t.box(/*#-1#*/-/*#-1#*/1, /*#-1#*/$g.____testlib.basictypes.Integer));
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      };
      $static.$equals = function (left, right) {
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#Boolean(NativeBoolean(left) == NativeBoolean(right))#*/$t.box(/*#NativeBoolean(left) == NativeBoolean(right))#*/$t.unbox(/*#left) == NativeBoolean(right))#*/left) == /*#NativeBoolean(right))#*/$t.unbox(/*#right))#*/right), /*#Boolean(NativeBoolean(left) == NativeBoolean(right))#*/$g.____testlib.basictypes.Boolean));
          return;
        };
        return $promise.new($continue);
      };
      $instance.String = function () {
        var $this = this;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#String(NativeBoolean(this).toString())#*/$t.box(/*#NativeBoolean(this).toString())#*/$t.unbox(/*#this).toString())#*/$this).toString(), /*#String(NativeBoolean(this).toString())#*/$g.____testlib.basictypes.String));
          return;
        };
        return $promise.new($continue);
      };
      $instance.MapKey = $t.property(function () {
        var $this = this;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#this }#*/$this);
          return;
        };
        return $promise.new($continue);
      });
      this.$typesig = function () {
        var computed = $t.createtypesig(['compare', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Integer).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['MapKey', 3, $g.____testlib.basictypes.Stringable.$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    this.$type('String', false, 'string', function () {
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
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#this#*/$this);
          return;
        };
        return $promise.new($continue);
      };
      $static.$equals = function (first, second) {
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#Boolean(NativeString(first) == NativeString(second))#*/$t.box(/*#NativeString(first) == NativeString(second))#*/$t.unbox(/*#first) == NativeString(second))#*/first) == /*#NativeString(second))#*/$t.unbox(/*#second))#*/second), /*#Boolean(NativeString(first) == NativeString(second))#*/$g.____testlib.basictypes.Boolean));
          return;
        };
        return $promise.new($continue);
      };
      $static.$plus = function (first, second) {
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#String(NativeString(first) + NativeString(second))#*/$t.box(/*#NativeString(first) + NativeString(second))#*/$t.unbox(/*#first) + NativeString(second))#*/first) + /*#NativeString(second))#*/$t.unbox(/*#second))#*/second), /*#String(NativeString(first) + NativeString(second))#*/$g.____testlib.basictypes.String));
          return;
        };
        return $promise.new($continue);
      };
      $instance.MapKey = $t.property(function () {
        var $this = this;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#this }#*/$this);
          return;
        };
        return $promise.new($continue);
      });
      $instance.Length = $t.property(function () {
        var $this = this;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          $resolve(/*#Integer(NativeString(this).length)#*/$t.box(/*#NativeString(this).length)#*/$t.unbox(/*#this).length)#*/$this).length, /*#Integer(NativeString(this).length)#*/$g.____testlib.basictypes.Integer));
          return;
        };
        return $promise.new($continue);
      });
      this.$typesig = function () {
        var computed = $t.createtypesig(['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['plus', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['MapKey', 3, $g.____testlib.basictypes.Stringable.$typeref()], ['Length', 3, $g.____testlib.basictypes.Integer.$typeref()]);
        this.$typesig = function () {
          return computed;
        };
        return computed;
      };
    });

    $static.formatTemplateString = function (pieces, values) {
      var $result;
      var $temp0;
      var $temp1;
      var i;
      var result;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              result = /*#''#*/$t.box(/*#''#*/'', /*#''#*/$g.____testlib.basictypes.String);
              $current = 1;
              continue;

            case 1:
/*#pieces.Length - 1) {#*/              pieces.Length().then(/*#pieces.Length - 1) {#*/function (/*#pieces.Length - 1) {#*/$result2) /*#pieces.Length - 1) {#*/{
                return /*#pieces.Length - 1) {#*/$g.____testlib.basictypes.Integer.$minus(/*#pieces.Length - 1) {#*/$result2, /*#1) {#*/$t.box(/*#1) {#*/1, /*#1) {#*/$g.____testlib.basictypes.Integer)).then(/*#1) {#*/function (/*#1) {#*/$result1) /*#1) {#*/{
                  return /*#0 .. (pieces.Length - 1) {#*/$g.____testlib.basictypes.Integer.$range(/*#0 .. (pieces.Length - 1) {#*/$t.box(/*#0 .. (pieces.Length - 1) {#*/0, /*#0 .. (pieces.Length - 1) {#*/$g.____testlib.basictypes.Integer), /*#pieces.Length - 1) {#*/$result1).then(/*#pieces.Length - 1) {#*/function (/*#pieces.Length - 1) {#*/$result0) /*#pieces.Length - 1) {#*/{
                    $result = /*#0 .. (pieces.Length - 1) {#*/$result0;
                    $current = 2;
                    $continue($resolve, $reject);
                    return;
                  });
                });
              }).catch(function (err) {
                $reject(err);
                return;
              });
              return;

            case 2:
              $temp1 = $result;
              $current = 3;
              continue;

            case 3:
/*#for i in 0 .. (pieces.Length - 1) {#*/              $temp1.Next().then(/*#for i in 0 .. (pieces.Length - 1) {#*/function (/*#for i in 0 .. (pieces.Length - 1) {#*/$result0) /*#for i in 0 .. (pieces.Length - 1) {#*/{
                $temp0 = /*#i in 0 .. (pieces.Length - 1) {#*/$result0;
                $result = /*#i in 0 .. (pieces.Length - 1) {#*/$temp0;
                $current = 4;
                $continue($resolve, $reject);
                return;
              }).catch(function (err) {
                $reject(err);
                return;
              });
              return;

            case 4:
/*#i in 0 .. (pieces.Length - 1) {#*/              i = /*#i in 0 .. (pieces.Length - 1) {#*/$temp0.First;
              if (/*#for i in 0 .. (pieces.Length - 1) {#*/$t.unbox(/*#for i in 0 .. (pieces.Length - 1) {#*/$temp0.Second)) /*#for i in 0 .. (pieces.Length - 1) {#*/{
                $current = 5;
                continue;
              } else {
                $current = 11;
                continue;
              }
              break;

            case 5:
/*#pieces[i]#*/              pieces.$index(/*#i]#*/i).then(/*#i]#*/function (/*#i]#*/$result1) /*#i]#*/{
                return /*#result + pieces[i]#*/$g.____testlib.basictypes.String.$plus(/*#result + pieces[i]#*/result, /*#pieces[i]#*/$result1).then(/*#pieces[i]#*/function (/*#pieces[i]#*/$result0) /*#pieces[i]#*/{
                  result = /*#result + pieces[i]#*/$result0;
                  $result = /*#result = result + pieces[i]#*/result;
                  $current = 6;
                  $continue($resolve, $reject);
                  return;
                });
              }).catch(function (err) {
                $reject(err);
                return;
              });
              return;

            case 6:
/*#values.Length {#*/              values.Length().then(/*#values.Length {#*/function (/*#values.Length {#*/$result1) /*#values.Length {#*/{
                return /*#i < values.Length {#*/$g.____testlib.basictypes.Integer.$compare(/*#i < values.Length {#*/i, /*#values.Length {#*/$result1).then(/*#values.Length {#*/function (/*#values.Length {#*/$result0) /*#values.Length {#*/{
                  $result = /*#i < values.Length {#*/$t.unbox(/*#i < values.Length {#*/$result0) < /*#i < values.Length {#*/0;
                  $current = 7;
                  $continue($resolve, $reject);
                  return;
                });
              }).catch(function (err) {
                $reject(err);
                return;
              });
              return;

            case 7:
              if ($result) {
                $current = 8;
                continue;
              } else {
                $current = 10;
                continue;
              }
              break;

            case 8:
/*#values[i].String()#*/              values.$index(/*#i].String()#*/i).then(/*#i].String()#*/function (/*#i].String()#*/$result2) /*#i].String()#*/{
                return /*#values[i].String()#*/$result2.String().then(/*#values[i].String()#*/function (/*#values[i].String()#*/$result1) /*#values[i].String()#*/{
                  return /*#result + values[i].String()#*/$g.____testlib.basictypes.String.$plus(/*#result + values[i].String()#*/result, /*#.String()#*/$result1).then(/*#.String()#*/function (/*#.String()#*/$result0) /*#.String()#*/{
                    result = /*#result + values[i].String()#*/$result0;
                    $result = /*#result = result + values[i].String()#*/result;
                    $current = 9;
                    $continue($resolve, $reject);
                    return;
                  });
                });
              }).catch(function (err) {
                $reject(err);
                return;
              });
              return;

            case 9:
              $current = 10;
              continue;

            case 10:
              $current = 3;
              continue;

            case 11:
              $resolve(/*#result#*/result);
              return;

            default:
              $resolve();
              return;
          }
        }
      };
      return $promise.new($continue);
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
                continue;

              case 1:
                $temp1 = /*#stream {#*/stream;
                $current = 2;
                continue;

              case 2:
/*#for item in stream {#*/                $temp1.Next().then(/*#for item in stream {#*/function (/*#for item in stream {#*/$result0) /*#for item in stream {#*/{
                  $temp0 = /*#item in stream {#*/$result0;
                  $result = /*#item in stream {#*/$temp0;
                  $current = 3;
                  $continue($yield, $yieldin, $reject, $done);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 3:
/*#item in stream {#*/                item = /*#item in stream {#*/$temp0.First;
                if (/*#for item in stream {#*/$t.unbox(/*#for item in stream {#*/$temp0.Second)) /*#for item in stream {#*/{
                  $current = 4;
                  continue;
                } else {
                  $current = 7;
                  continue;
                }
                break;

              case 4:
/*#mapper(item)#*/                mapper(/*#item)#*/item).then(/*#item)#*/function (/*#item)#*/$result0) /*#item)#*/{
                  $result = /*#mapper(item)#*/$result0;
                  $current = 5;
                  $continue($yield, $yieldin, $reject, $done);
                  return;
                }).catch(function (err) {
                  $reject(err);
                  return;
                });
                return;

              case 5:
                $yield($result);
                $current = 6;
                return;

              case 6:
                $current = 2;
                continue;

              default:
                $done();
                return;
            }
          }
        };
        return $generator.new($continue);
      };
      return $f;
    };
  });
  $module('basic', function () {
    var $static = this;
    $static.DoSomething = function () {
      var $result;
      var bar;
      var foo;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              foo = /*#1#*/$t.box(/*#1#*/1, /*#1#*/$g.____testlib.basictypes.Integer);
              bar = /*#'hi there!'#*/$t.box(/*#'hi there!'#*/'hi there!', /*#'hi there!'#*/$g.____testlib.basictypes.String);
/*#foo == 2 && bar == 'hello world'#*/              $g.____testlib.basictypes.Integer.$equals(/*#foo == 2 && bar == 'hello world'#*/foo, /*#2 && bar == 'hello world'#*/$t.box(/*#2 && bar == 'hello world'#*/2, /*#2 && bar == 'hello world'#*/$g.____testlib.basictypes.Integer)).then(/*#2 && bar == 'hello world'#*/function (/*#2 && bar == 'hello world'#*/$result1) /*#2 && bar == 'hello world'#*/{
                return /*#foo == 2 && bar == 'hello world'#*/$promise.resolve(/*#foo == 2 && bar == 'hello world'#*/$t.unbox(/*#foo == 2 && bar == 'hello world'#*/$result1)).then(/*#foo == 2 && bar == 'hello world'#*/function (/*#foo == 2 && bar == 'hello world'#*/$result0) /*#foo == 2 && bar == 'hello world'#*/{
                  return /*#foo == 2 && bar == 'hello world'#*/(/*#foo == 2 && bar == 'hello world'#*/$promise.shortcircuit(/*#foo == 2 && bar == 'hello world'#*/$result0, /*#foo == 2 && bar == 'hello world'#*/true) || /*#bar == 'hello world'#*/$g.____testlib.basictypes.String.$equals(/*#bar == 'hello world'#*/bar, /*#'hello world'#*/$t.box(/*#'hello world'#*/'hello world', /*#'hello world'#*/$g.____testlib.basictypes.String))).then(/*#'hello world'#*/function (/*#'hello world'#*/$result2) /*#'hello world'#*/{
                    $result = /*#foo == 2 && bar == 'hello world'#*/$t.box(/*#foo == 2 && bar == 'hello world'#*/$result0 && /*#foo == 2 && bar == 'hello world'#*/$t.unbox(/*#bar == 'hello world'#*/$result2), /*#foo == 2 && bar == 'hello world'#*/$g.____testlib.basictypes.Boolean);
                    $current = 1;
                    $continue($resolve, $reject);
                    return;
                  });
                });
              }).catch(function (err) {
                $reject(err);
                return;
              });
              return;

            case 1:
              $resolve($result);
              return;

            default:
              $resolve();
              return;
          }
        }
      };
      return $promise.new($continue);
    };
  });
  $module('____testlib.basic.webidl', function () {
    var $static = this;







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
        $global.close();
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

