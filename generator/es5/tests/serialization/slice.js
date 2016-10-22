$module('slice', function () {
  var $static = this;
  this.$struct('9b59dab3', 'AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherInt) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        AnotherInt: AnotherInt,
      };
      instance.$markruntimecreated();
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'AnotherInt', 'AnotherInt', function () {
      return $g.____testlib.basictypes.Integer;
    }, true, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<9b59dab3>": true,
        "equals|4|29dc432d<5ab5941e>": true,
        "Stringify|2|29dc432d<538656f2>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<9b59dab3>": true,
        "String|2|29dc432d<538656f2>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('a7573ae2', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Values) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        Values: Values,
      };
      instance.$markruntimecreated();
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'Values', 'Values', function () {
      return $g.____testlib.basictypes.Slice($g.slice.AnotherStruct);
    }, true, function () {
      return $g.____testlib.basictypes.Slice($g.slice.AnotherStruct);
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<a7573ae2>": true,
        "equals|4|29dc432d<5ab5941e>": true,
        "Stringify|2|29dc432d<538656f2>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<a7573ae2>": true,
        "String|2|29dc432d<538656f2>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $result;
    var correct;
    var jsonString;
    var parsed;
    var s;
    var values;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.slice.AnotherStruct.new($t.box(1, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $g.slice.AnotherStruct.new($t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result2) {
                return $g.slice.AnotherStruct.new($t.box(3, $g.____testlib.basictypes.Integer)).then(function ($result3) {
                  return $g.____testlib.basictypes.List($g.slice.AnotherStruct).forArray([$result1, $result2, $result3]).then(function ($result0) {
                    $result = $result0;
                    $current = 1;
                    $continue($resolve, $reject);
                    return;
                  });
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            values = $result;
            values.$slice($t.box(0, $g.____testlib.basictypes.Integer), null).then(function ($result1) {
              return $g.slice.SomeStruct.new($result1).then(function ($result0) {
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
            s = $result;
            jsonString = $t.box('{"Values":[{"AnotherInt":1},{"AnotherInt":2},{"AnotherInt":3}]}', $g.____testlib.basictypes.String);
            s.Stringify($g.____testlib.basictypes.JSON)().then(function ($result1) {
              return $g.____testlib.basictypes.String.$equals($result1, jsonString).then(function ($result0) {
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
            correct = $result;
            $g.slice.SomeStruct.Parse($g.____testlib.basictypes.JSON)(jsonString).then(function ($result0) {
              $result = $result0;
              $current = 4;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 4:
            parsed = $result;
            $promise.resolve($t.unbox(correct)).then(function ($result1) {
              return ($promise.shortcircuit($result1, true) || s.Values.Length()).then(function ($result3) {
                return ($promise.shortcircuit($result1, true) || $g.____testlib.basictypes.Integer.$equals($result3, $t.box(3, $g.____testlib.basictypes.Integer))).then(function ($result2) {
                  return $promise.resolve($result1 && $t.unbox($result2)).then(function ($result0) {
                    return ($promise.shortcircuit($result0, true) || s.Values.$index($t.box(0, $g.____testlib.basictypes.Integer))).then(function ($result5) {
                      return ($promise.shortcircuit($result0, true) || $g.____testlib.basictypes.Integer.$equals($result5.AnotherInt, $t.box(1, $g.____testlib.basictypes.Integer))).then(function ($result4) {
                        $result = $t.box($result0 && $t.unbox($result4), $g.____testlib.basictypes.Boolean);
                        $current = 5;
                        $continue($resolve, $reject);
                        return;
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

          case 5:
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
