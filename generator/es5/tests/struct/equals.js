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
    this.$typesig = function () {
      return $t.createtypesig(['SomeValue', 5, $g.____testlib.basictypes.Integer.$typeref()], ['AnotherValue', 5, $g.equals.Bar.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.equals.Foo).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.equals.Foo).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
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
    this.$typesig = function () {
      return $t.createtypesig(['StringValue', 5, $g.____testlib.basictypes.String.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.equals.Bar).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.equals.Bar).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
  });

  $static.TEST = function () {
    var copy;
    var different;
    var first;
    var second;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.equals.Bar.new($t.box('hello world', $g.____testlib.basictypes.String)).then(function ($result1) {
              $temp0 = $result1;
              return $g.equals.Foo.new($t.box(42, $g.____testlib.basictypes.Integer), ($temp0, $temp0)).then(function ($result0) {
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
            first = $result;
            second = first;
            $g.equals.Bar.new($t.box('hello world', $g.____testlib.basictypes.String)).then(function ($result1) {
              $temp2 = $result1;
              return $g.equals.Foo.new($t.box(42, $g.____testlib.basictypes.Integer), ($temp2, $temp2)).then(function ($result0) {
                $temp3 = $result0;
                $result = ($temp3, $temp3);
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
            copy = $result;
            $g.equals.Bar.new($t.box('hello worlds!', $g.____testlib.basictypes.String)).then(function ($result1) {
              $temp4 = $result1;
              return $g.equals.Foo.new($t.box(42, $g.____testlib.basictypes.Integer), ($temp4, $temp4)).then(function ($result0) {
                $temp5 = $result0;
                $result = ($temp5, $temp5);
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
            different = $result;
            $g.equals.Foo.$equals(first, second).then(function ($result3) {
              return $promise.resolve($t.unbox($result3)).then(function ($result2) {
                return ($promise.shortcircuit(!$result2) || $g.equals.Foo.$equals(first, copy)).then(function ($result4) {
                  return $promise.resolve($result2 && $t.unbox($result4)).then(function ($result1) {
                    return ($promise.shortcircuit(!$result1) || $g.equals.Foo.$equals(first, different)).then(function ($result5) {
                      return $promise.resolve($result1 && !$t.unbox($result5)).then(function ($result0) {
                        return ($promise.shortcircuit(!$result0) || $g.equals.Foo.$equals(copy, different)).then(function ($result6) {
                          $result = $t.box($result0 && !$t.unbox($result6), $g.____testlib.basictypes.Boolean);
                          $current = 4;
                          $continue($resolve, $reject);
                          return;
                        });
                      });
                    });
                  });
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 4:
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
