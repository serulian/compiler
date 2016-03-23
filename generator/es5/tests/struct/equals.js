$module('equals', function () {
  var $static = this;
  this.$struct('Foo', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeValue, AnotherValue) {
      var instance = new $static();
      instance.$data = {
      };
      instance.$lazycheck = false;
      instance.SomeValue = SomeValue;
      instance.AnotherValue = AnotherValue;
      return $promise.resolve(instance);
    };
    $instance.Mapping = function () {
      var mappedData = {
      };
      mappedData['SomeValue'] = this.SomeValue;
      mappedData['AnotherValue'] = this.AnotherValue;
      return $promise.resolve($t.nominalwrap(mappedData, $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left.$data['SomeValue'], right.$data['SomeValue'], $g.____testlib.basictypes.Integer));
      promises.push($t.equals(left.$data['AnotherValue'], right.$data['AnotherValue'], $g.equals.Bar));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
          }
        }
        return $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'SomeValue', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['SomeValue'], $g.____testlib.basictypes.Integer, false, 'SomeValue');
        }
        if (this.$data['SomeValue'] != null) {
          return $t.box(this.$data['SomeValue'], $g.____testlib.basictypes.Integer);
        }
        return this.$data['SomeValue'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['SomeValue'] = $t.unbox(val);
          return;
        }
        this.$data['SomeValue'] = val;
      },
    });
    Object.defineProperty($instance, 'AnotherValue', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['AnotherValue'], $g.equals.Bar, false, 'AnotherValue');
        }
        if (this.$data['AnotherValue'] != null) {
          return $t.box(this.$data['AnotherValue'], $g.equals.Bar);
        }
        return this.$data['AnotherValue'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['AnotherValue'] = $t.unbox(val);
          return;
        }
        this.$data['AnotherValue'] = val;
      },
    });
  });

  this.$struct('Bar', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (StringValue) {
      var instance = new $static();
      instance.$data = {
      };
      instance.$lazycheck = false;
      instance.StringValue = StringValue;
      return $promise.resolve(instance);
    };
    $instance.Mapping = function () {
      var mappedData = {
      };
      mappedData['StringValue'] = this.StringValue;
      return $promise.resolve($t.nominalwrap(mappedData, $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left.$data['StringValue'], right.$data['StringValue'], $g.____testlib.basictypes.String));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
          }
        }
        return $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'StringValue', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['StringValue'], $g.____testlib.basictypes.String, false, 'StringValue');
        }
        if (this.$data['StringValue'] != null) {
          return $t.box(this.$data['StringValue'], $g.____testlib.basictypes.String);
        }
        return this.$data['StringValue'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['StringValue'] = $t.unbox(val);
          return;
        }
        this.$data['StringValue'] = val;
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
            $g.equals.Bar.new($t.nominalwrap('hello world', $g.____testlib.basictypes.String)).then(function ($result0) {
              $temp0 = $result0;
              return $g.equals.Foo.new($t.nominalwrap(42, $g.____testlib.basictypes.Integer), ($temp0, $temp0)).then(function ($result1) {
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
            $g.equals.Bar.new($t.nominalwrap('hello world', $g.____testlib.basictypes.String)).then(function ($result0) {
              $temp2 = $result0;
              return $g.equals.Foo.new($t.nominalwrap(42, $g.____testlib.basictypes.Integer), ($temp2, $temp2)).then(function ($result1) {
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
            $g.equals.Bar.new($t.nominalwrap('hello worlds!', $g.____testlib.basictypes.String)).then(function ($result0) {
              $temp4 = $result0;
              return $g.equals.Foo.new($t.nominalwrap(42, $g.____testlib.basictypes.Integer), ($temp4, $temp4)).then(function ($result1) {
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
            $g.equals.Foo.$equals(first, second).then(function ($result0) {
              return $g.equals.Foo.$equals(first, copy).then(function ($result1) {
                return $g.equals.Foo.$equals(first, different).then(function ($result2) {
                  return $g.equals.Foo.$equals(copy, different).then(function ($result3) {
                    $result = $t.nominalwrap((($t.nominalunwrap($result0) && $t.nominalunwrap($result1)) && !$t.nominalunwrap($result2)) && !$t.nominalunwrap($result3), $g.____testlib.basictypes.Boolean);
                    $state.current = 4;
                    $callback($state);
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
