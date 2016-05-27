$module('basic', function () {
  var $static = this;
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
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['SomeInstance'], $g.basic.AnotherStruct, false, 'SomeInstance');
            this.$lazychecked['SomeInstance'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['SomeInstance'], $g.basic.AnotherStruct);
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
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['SomeInstance'], right[BOXED_DATA_PROPERTY]['SomeInstance'], $g.basic.AnotherStruct));
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
    var ss;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.basic.AnotherStruct.new($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result1) {
              $temp0 = $result1;
              return $g.basic.SomeStruct.new($t.box(42, $g.____testlib.basictypes.Integer), $t.box(true, $g.____testlib.basictypes.Boolean), ($temp0, $temp0)).then(function ($result0) {
                $temp1 = $result0;
                $result = ($temp1, $temp1);
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
            ss = $result;
            $g.____testlib.basictypes.Integer.$equals(ss.SomeField, $t.box(42, $g.____testlib.basictypes.Integer)).then(function ($result2) {
              return $promise.resolve($t.unbox($result2)).then(function ($result1) {
                return $promise.resolve($result1 && $t.unbox(ss.AnotherField)).then(function ($result0) {
                  $result = $t.box($result0 && $t.unbox(ss.SomeInstance.AnotherBool), $g.____testlib.basictypes.Boolean);
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
