$module('custom', function () {
  var $static = this;
  this.$class('CustomJSON', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.Get = function () {
      var $state = $t.sm(function ($continue) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.custom.CustomJSON.new().then(function ($result0) {
                $result = $result0;
                $state.current = 1;
                $continue($state);
              }).catch(function (err) {
                $state.reject(err);
              });
              return;

            case 1:
              $state.resolve($result);
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    };
    $instance.Stringify = function (value) {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.____testlib.basictypes.JSON.Get().then(function ($result0) {
                return $result0.Stringify(value).then(function ($result1) {
                  $result = $result1;
                  $state.current = 1;
                  $continue($state);
                });
              }).catch(function (err) {
                $state.reject(err);
              });
              return;

            case 1:
              $state.resolve($result);
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    };
    $instance.Parse = function (value) {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.____testlib.basictypes.JSON.Get().then(function ($result0) {
                return $result0.Parse(value).then(function ($result1) {
                  $result = $result1;
                  $state.current = 1;
                  $continue($state);
                });
              }).catch(function (err) {
                $state.reject(err);
              });
              return;

            case 1:
              $state.resolve($result);
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    };
  });

  this.$struct('AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherBool) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        AnotherBool: AnotherBool,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'AnotherBool', {
        get: function () {
          if (this.$lazychecked['AnotherBool']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['AnotherBool'], $g.____testlib.basictypes.Boolean, false, 'AnotherBool');
            this.$lazychecked['AnotherBool'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['AnotherBool'], $g.____testlib.basictypes.Boolean);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['AnotherBool'] = this.AnotherBool;
        return $promise.resolve($t.box(mapped, $g.____testlib.basictypes.Mapping($t.any)));
      };
      return instance;
    };
    $instance.Mapping = function () {
      return $promise.resolve($t.box(this[BOXED_DATA_PROPERTY], $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['AnotherBool'], right[BOXED_DATA_PROPERTY]['AnotherBool'], $g.____testlib.basictypes.Boolean));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'AnotherBool', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['AnotherBool'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['AnotherBool'] = value;
      },
    });
  });

  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField, AnotherField, SomeInstance) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        SomeField: SomeField,
        AnotherField: AnotherField,
        SomeInstance: SomeInstance,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'SomeField', {
        get: function () {
          if (this.$lazychecked['SomeField']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['SomeField'], $g.____testlib.basictypes.Integer, false, 'SomeField');
            this.$lazychecked['SomeField'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['SomeField'], $g.____testlib.basictypes.Integer);
        },
      });
      Object.defineProperty(instance, 'AnotherField', {
        get: function () {
          if (this.$lazychecked['AnotherField']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['AnotherField'], $g.____testlib.basictypes.Boolean, false, 'AnotherField');
            this.$lazychecked['AnotherField'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['AnotherField'], $g.____testlib.basictypes.Boolean);
        },
      });
      Object.defineProperty(instance, 'SomeInstance', {
        get: function () {
          if (this.$lazychecked['SomeInstance']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['SomeInstance'], $g.custom.AnotherStruct, false, 'SomeInstance');
            this.$lazychecked['SomeInstance'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['SomeInstance'], $g.custom.AnotherStruct);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['SomeField'] = this.SomeField;
        mapped['AnotherField'] = this.AnotherField;
        mapped['SomeInstance'] = this.SomeInstance;
        return $promise.resolve($t.box(mapped, $g.____testlib.basictypes.Mapping($t.any)));
      };
      return instance;
    };
    $instance.Mapping = function () {
      return $promise.resolve($t.box(this[BOXED_DATA_PROPERTY], $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['SomeField'], right[BOXED_DATA_PROPERTY]['SomeField'], $g.____testlib.basictypes.Integer));
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['AnotherField'], right[BOXED_DATA_PROPERTY]['AnotherField'], $g.____testlib.basictypes.Boolean));
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['SomeInstance'], right[BOXED_DATA_PROPERTY]['SomeInstance'], $g.custom.AnotherStruct));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'SomeField', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['SomeField'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['SomeField'] = value;
      },
    });
    Object.defineProperty($instance, 'AnotherField', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['AnotherField'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['AnotherField'] = value;
      },
    });
    Object.defineProperty($instance, 'SomeInstance', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['SomeInstance'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['SomeInstance'] = value;
      },
    });
  });

  $static.TEST = function () {
    var jsonString;
    var parsed;
    var s;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.custom.AnotherStruct.new($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $temp0 = $result0;
              return $g.custom.SomeStruct.new($t.box(2, $g.____testlib.basictypes.Integer), $t.box(false, $g.____testlib.basictypes.Boolean), ($temp0, $temp0)).then(function ($result1) {
                $temp1 = $result1;
                $result = ($temp1, $temp1);
                $state.current = 1;
                $continue($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            s = $result;
            jsonString = $t.box('{"AnotherField":false,"SomeField":2,"SomeInstance":{"AnotherBool":true}}', $g.____testlib.basictypes.String);
            $g.custom.SomeStruct.Parse($g.custom.CustomJSON)(jsonString).then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            parsed = $result;
            $g.____testlib.basictypes.Integer.$equals(parsed.SomeField, $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result2) {
              return $promise.resolve($t.unbox($result2)).then(function ($result1) {
                return $promise.resolve($result1 && !$t.unbox(parsed.AnotherField)).then(function ($result0) {
                  $result = $t.box($result0 && $t.unbox(parsed.SomeInstance.AnotherBool), $g.____testlib.basictypes.Boolean);
                  $state.current = 3;
                  $continue($state);
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            $state.resolve($result);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
