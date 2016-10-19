$module('equals', function () {
  var $static = this;
  this.$struct('Foo', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeValue, AnotherValue) {
      var instance = new $static();
      instance.$unboxed = false;
      instance[BOXED_DATA_PROPERTY] = {
        SomeValue: SomeValue,
        AnotherValue: AnotherValue,
      };
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeValue', 'SomeValue', function () {
      return $g.____testlib.basictypes.Integer;
    }, true, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'AnotherValue', 'AnotherValue', function () {
      return $g.equals.Bar;
    }, true, function () {
      return $g.equals.Bar;
    }, false);
    this.$typesig = function () {
      var computed = $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.equals.Foo).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.equals.Foo).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.equals.Foo).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
    };
  });

  this.$struct('Bar', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (StringValue) {
      var instance = new $static();
      instance.$unboxed = false;
      instance[BOXED_DATA_PROPERTY] = {
        StringValue: StringValue,
      };
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'StringValue', 'StringValue', function () {
      return $g.____testlib.basictypes.String;
    }, true, function () {
      return $g.____testlib.basictypes.String;
    }, false);
    this.$typesig = function () {
      var computed = $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.equals.Bar).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.equals.Bar).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.equals.Bar).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
    };
  });

  $static.TEST = function () {
    var $result;
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
              return $g.equals.Foo.new($t.box(42, $g.____testlib.basictypes.Integer), $result1).then(function ($result0) {
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
            first = $result;
            second = first;
            $g.equals.Bar.new($t.box('hello world', $g.____testlib.basictypes.String)).then(function ($result1) {
              return $g.equals.Foo.new($t.box(42, $g.____testlib.basictypes.Integer), $result1).then(function ($result0) {
                $result = $result0;
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
              return $g.equals.Foo.new($t.box(42, $g.____testlib.basictypes.Integer), $result1).then(function ($result0) {
                $result = $result0;
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
                return ($promise.shortcircuit($result2, true) || $g.equals.Foo.$equals(first, copy)).then(function ($result4) {
                  return $promise.resolve($result2 && $t.unbox($result4)).then(function ($result1) {
                    return ($promise.shortcircuit($result1, true) || $g.equals.Foo.$equals(first, different)).then(function ($result5) {
                      return $promise.resolve($result1 && !$t.unbox($result5)).then(function ($result0) {
                        return ($promise.shortcircuit($result0, true) || $g.equals.Foo.$equals(copy, different)).then(function ($result6) {
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
