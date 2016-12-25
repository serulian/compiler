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
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'AnotherInt', 'AnotherInt', function () {
      return $g.____testlib.basictypes.Integer;
    }, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<9b59dab3>": true,
        "equals|4|29dc432d<43834c3f>": true,
        "Stringify|2|29dc432d<5cffd9b5>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<9b59dab3>": true,
        "String|2|29dc432d<5cffd9b5>": true,
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
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'Values', 'Values', function () {
      return $g.____testlib.basictypes.Slice($g.slice.AnotherStruct);
    }, function () {
      return $g.____testlib.basictypes.Slice($g.slice.AnotherStruct);
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<a7573ae2>": true,
        "equals|4|29dc432d<43834c3f>": true,
        "Stringify|2|29dc432d<5cffd9b5>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<a7573ae2>": true,
        "String|2|29dc432d<5cffd9b5>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = $t.markpromising(function () {
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
            values = $g.____testlib.basictypes.List($g.slice.AnotherStruct).forArray([$g.slice.AnotherStruct.new($t.fastbox(1, $g.____testlib.basictypes.Integer)), $g.slice.AnotherStruct.new($t.fastbox(2, $g.____testlib.basictypes.Integer)), $g.slice.AnotherStruct.new($t.fastbox(3, $g.____testlib.basictypes.Integer))]);
            s = $g.slice.SomeStruct.new(values.$slice($t.fastbox(0, $g.____testlib.basictypes.Integer), null));
            jsonString = $t.fastbox('{"Values":[{"AnotherInt":1},{"AnotherInt":2},{"AnotherInt":3}]}', $g.____testlib.basictypes.String);
            $promise.maybe(s.Stringify($g.____testlib.basictypes.JSON)()).then(function ($result0) {
              $result = $g.____testlib.basictypes.String.$equals($result0, jsonString);
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            correct = $result;
            $promise.maybe($g.slice.SomeStruct.Parse($g.____testlib.basictypes.JSON)(jsonString)).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            parsed = $result;
            $resolve($t.fastbox((correct.$wrapped && (s.Values.Length().$wrapped == 3)) && (s.Values.$index($t.fastbox(0, $g.____testlib.basictypes.Integer)).AnotherInt.$wrapped == 1), $g.____testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
