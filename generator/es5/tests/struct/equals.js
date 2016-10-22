$module('equals', function () {
  var $static = this;
  this.$struct('9285bf36', 'Foo', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeValue, AnotherValue) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        SomeValue: SomeValue,
        AnotherValue: AnotherValue,
      };
      instance.$markruntimecreated();
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeValue', 'SomeValue', function () {
      return $g.____testlib.basictypes.Integer;
    }, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'AnotherValue', 'AnotherValue', function () {
      return $g.equals.Bar;
    }, function () {
      return $g.equals.Bar;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<9285bf36>": true,
        "equals|4|29dc432d<5ab5941e>": true,
        "Stringify|2|29dc432d<538656f2>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<9285bf36>": true,
        "String|2|29dc432d<538656f2>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('4fb96a52', 'Bar', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (StringValue) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        StringValue: StringValue,
      };
      instance.$markruntimecreated();
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'StringValue', 'StringValue', function () {
      return $g.____testlib.basictypes.String;
    }, function () {
      return $g.____testlib.basictypes.String;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<4fb96a52>": true,
        "equals|4|29dc432d<5ab5941e>": true,
        "Stringify|2|29dc432d<538656f2>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<4fb96a52>": true,
        "String|2|29dc432d<538656f2>": true,
      };
      return this.$cachedtypesig = computed;
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
            $g.equals.Bar.new($t.fastbox('hello world', $g.____testlib.basictypes.String)).then(function ($result1) {
              return $g.equals.Foo.new($t.fastbox(42, $g.____testlib.basictypes.Integer), $result1).then(function ($result0) {
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
            $g.equals.Bar.new($t.fastbox('hello world', $g.____testlib.basictypes.String)).then(function ($result1) {
              return $g.equals.Foo.new($t.fastbox(42, $g.____testlib.basictypes.Integer), $result1).then(function ($result0) {
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
            $g.equals.Bar.new($t.fastbox('hello worlds!', $g.____testlib.basictypes.String)).then(function ($result1) {
              return $g.equals.Foo.new($t.fastbox(42, $g.____testlib.basictypes.Integer), $result1).then(function ($result0) {
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
              return $promise.resolve($result3.$wrapped).then(function ($result2) {
                return ($promise.shortcircuit($result2, true) || $g.equals.Foo.$equals(first, copy)).then(function ($result4) {
                  return $promise.resolve($result2 && $result4.$wrapped).then(function ($result1) {
                    return ($promise.shortcircuit($result1, true) || $g.equals.Foo.$equals(first, different)).then(function ($result5) {
                      return $promise.resolve($result1 && !$result5.$wrapped).then(function ($result0) {
                        return ($promise.shortcircuit($result0, true) || $g.equals.Foo.$equals(copy, different)).then(function ($result6) {
                          $result = $t.fastbox($result0 && !$result6.$wrapped, $g.____testlib.basictypes.Boolean);
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
