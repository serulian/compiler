$module('structprops', function () {
  var $static = this;
  this.$struct('SomeProps', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (BoolValue, StringValue) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        BoolValue: BoolValue,
        StringValue: StringValue,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'BoolValue', {
        get: function () {
          if (this.$lazychecked['BoolValue']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['BoolValue'], $g.____testlib.basictypes.Boolean, false, 'BoolValue');
            this.$lazychecked['BoolValue'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['BoolValue'], $g.____testlib.basictypes.Boolean);
        },
      });
      Object.defineProperty(instance, 'StringValue', {
        get: function () {
          if (this.$lazychecked['StringValue']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['StringValue'], $g.____testlib.basictypes.String, false, 'StringValue');
            this.$lazychecked['StringValue'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['StringValue'], $g.____testlib.basictypes.String);
        },
      });
      Object.defineProperty(instance, 'OptionalValue', {
        get: function () {
          if (this.$lazychecked['OptionalValue']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['OptionalValue'], $g.____testlib.basictypes.Integer, true, 'OptionalValue');
            this.$lazychecked['OptionalValue'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['OptionalValue'], $g.____testlib.basictypes.Integer);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['BoolValue'] = this.BoolValue;
        mapped['StringValue'] = this.StringValue;
        mapped['OptionalValue'] = this.OptionalValue;
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
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['BoolValue'], right[BOXED_DATA_PROPERTY]['BoolValue'], $g.____testlib.basictypes.Boolean));
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['StringValue'], right[BOXED_DATA_PROPERTY]['StringValue'], $g.____testlib.basictypes.String));
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['OptionalValue'], right[BOXED_DATA_PROPERTY]['OptionalValue'], $g.____testlib.basictypes.Integer));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'BoolValue', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['BoolValue'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['BoolValue'] = value;
      },
    });
    Object.defineProperty($instance, 'StringValue', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['StringValue'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['StringValue'] = value;
      },
    });
    Object.defineProperty($instance, 'OptionalValue', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['OptionalValue'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['OptionalValue'] = value;
      },
    });
    this.$typesig = function () {
      return $t.createtypesig(['BoolValue', 5, $g.____testlib.basictypes.Boolean.$typeref()], ['StringValue', 5, $g.____testlib.basictypes.String.$typeref()], ['OptionalValue', 5, $g.____testlib.basictypes.Integer.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.structprops.SomeProps).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.structprops.SomeProps).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.structprops.SomeProps).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
  });

  $static.SimpleFunction = function (props) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.String.$equals(props.StringValue, $t.box("hello world", $g.____testlib.basictypes.String)).then(function ($result2) {
              return $promise.resolve($t.unbox($result2)).then(function ($result1) {
                return $promise.resolve($result1 && $t.unbox(props.BoolValue)).then(function ($result0) {
                  $result = $t.box($result0 && !(props.OptionalValue == null), $g.____testlib.basictypes.Boolean);
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
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.structprops.SomeProps.new(true, $t.box("hello world", $g.____testlib.basictypes.String)).then(function ($result1) {
              $temp0 = $result1;
              return $g.structprops.SimpleFunction(($temp0, $temp0.OptionalValue = $t.box(42, $g.____testlib.basictypes.Integer), $temp0)).then(function ($result0) {
                $result = $result0;
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
