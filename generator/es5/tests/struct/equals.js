$module('equals', function () {
  var $static = this;
  this.$struct('Foo', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeValue, AnotherValue) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        SomeValue: SomeValue,
        AnotherValue: AnotherValue,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'SomeValue', {
        get: function () {
          if (this.$lazychecked['SomeValue']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['SomeValue'], $g.____testlib.basictypes.Integer, false, 'SomeValue');
            this.$lazychecked['SomeValue'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['SomeValue'], $g.____testlib.basictypes.Integer);
        },
      });
      Object.defineProperty(instance, 'AnotherValue', {
        get: function () {
          if (this.$lazychecked['AnotherValue']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['AnotherValue'], $g.equals.Bar, false, 'AnotherValue');
            this.$lazychecked['AnotherValue'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['AnotherValue'], $g.equals.Bar);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['SomeValue'] = this.SomeValue;
        mapped['AnotherValue'] = this.AnotherValue;
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
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['SomeValue'], right[BOXED_DATA_PROPERTY]['SomeValue'], $g.____testlib.basictypes.Integer));
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['AnotherValue'], right[BOXED_DATA_PROPERTY]['AnotherValue'], $g.equals.Bar));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'SomeValue', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['SomeValue'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['SomeValue'] = value;
      },
    });
    Object.defineProperty($instance, 'AnotherValue', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['AnotherValue'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['AnotherValue'] = value;
      },
    });
  });

  this.$struct('Bar', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (StringValue) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        StringValue: StringValue,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'StringValue', {
        get: function () {
          if (this.$lazychecked['StringValue']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['StringValue'], $g.____testlib.basictypes.String, false, 'StringValue');
            this.$lazychecked['StringValue'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['StringValue'], $g.____testlib.basictypes.String);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['StringValue'] = this.StringValue;
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
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['StringValue'], right[BOXED_DATA_PROPERTY]['StringValue'], $g.____testlib.basictypes.String));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'StringValue', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['StringValue'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['StringValue'] = value;
      },
    });
  });

  $static.TEST = function () {
    var copy;
    var different;
    var first;
    var second;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.equals.Bar.new($t.box('hello world', $g.____testlib.basictypes.String)).then(function ($result0) {
              $temp0 = $result0;
              return $g.equals.Foo.new($t.box(42, $g.____testlib.basictypes.Integer), ($temp0, $temp0)).then(function ($result1) {
                $temp1 = $result1;
                $result = ($temp1, $temp1);
                $state.current = 1;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            first = $result;
            second = first;
            $g.equals.Bar.new($t.box('hello world', $g.____testlib.basictypes.String)).then(function ($result0) {
              $temp2 = $result0;
              return $g.equals.Foo.new($t.box(42, $g.____testlib.basictypes.Integer), ($temp2, $temp2)).then(function ($result1) {
                $temp3 = $result1;
                $result = ($temp3, $temp3);
                $state.current = 2;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            copy = $result;
            $g.equals.Bar.new($t.box('hello worlds!', $g.____testlib.basictypes.String)).then(function ($result0) {
              $temp4 = $result0;
              return $g.equals.Foo.new($t.box(42, $g.____testlib.basictypes.Integer), ($temp4, $temp4)).then(function ($result1) {
                $temp5 = $result1;
                $result = ($temp5, $temp5);
                $state.current = 3;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            different = $result;
            $g.equals.Foo.$equals(first, second).then(function ($result3) {
              return $promise.resolve($t.unbox($result3)).then(function ($result2) {
                return ($promise.shortcircuit($result2, false) || $g.equals.Foo.$equals(first, copy)).then(function ($result4) {
                  return $promise.resolve($result2 && $t.unbox($result4)).then(function ($result1) {
                    return ($promise.shortcircuit($result1, false) || $g.equals.Foo.$equals(first, different)).then(function ($result5) {
                      return $promise.resolve($result1 && !$t.unbox($result5)).then(function ($result0) {
                        return ($promise.shortcircuit($result0, false) || $g.equals.Foo.$equals(copy, different)).then(function ($result6) {
                          $result = $t.box($result0 && !$t.unbox($result6), $g.____testlib.basictypes.Boolean);
                          $state.current = 4;
                          $callback($state);
                        });
                      });
                    });
                  });
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 4:
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
