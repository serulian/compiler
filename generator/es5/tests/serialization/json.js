$module('json', function () {
  var $static = this;
  this.$struct('143c410b', 'AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherBool) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        AnotherBool: AnotherBool,
      };
      instance.$markruntimecreated();
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'AnotherBool', 'AnotherBool', function () {
      return $g.____testlib.basictypes.Boolean;
    }, function () {
      return $g.____testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<143c410b>": true,
        "equals|4|29dc432d<5ab5941e>": true,
        "Stringify|2|29dc432d<538656f2>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<143c410b>": true,
        "String|2|29dc432d<538656f2>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('8c772f6d', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField, AnotherField, SomeInstance) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        SomeField: SomeField,
        AnotherField: AnotherField,
        SomeInstance: SomeInstance,
      };
      instance.$markruntimecreated();
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'SomeField', function () {
      return $g.____testlib.basictypes.Integer;
    }, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'AnotherField', 'AnotherField', function () {
      return $g.____testlib.basictypes.Boolean;
    }, function () {
      return $g.____testlib.basictypes.Boolean;
    }, false);
    $t.defineStructField($static, 'SomeInstance', 'SomeInstance', function () {
      return $g.json.AnotherStruct;
    }, function () {
      return $g.json.AnotherStruct;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<8c772f6d>": true,
        "equals|4|29dc432d<5ab5941e>": true,
        "Stringify|2|29dc432d<538656f2>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<8c772f6d>": true,
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
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.json.AnotherStruct.new($t.fastbox(true, $g.____testlib.basictypes.Boolean)).then(function ($result1) {
              return $g.json.SomeStruct.new($t.fastbox(2, $g.____testlib.basictypes.Integer), $t.fastbox(false, $g.____testlib.basictypes.Boolean), $result1).then(function ($result0) {
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
            s = $result;
            jsonString = $t.fastbox('{"AnotherField":false,"SomeField":2,"SomeInstance":{"AnotherBool":true}}', $g.____testlib.basictypes.String);
            s.Stringify($g.____testlib.basictypes.JSON)().then(function ($result1) {
              return $g.____testlib.basictypes.String.$equals($result1, jsonString).then(function ($result0) {
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
            correct = $result;
            $g.json.SomeStruct.Parse($g.____testlib.basictypes.JSON)(jsonString).then(function ($result0) {
              $result = $result0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            parsed = $result;
            $promise.resolve(correct.$wrapped).then(function ($result2) {
              return ($promise.shortcircuit($result2, true) || $g.____testlib.basictypes.Integer.$equals(parsed.SomeField, $t.fastbox(2, $g.____testlib.basictypes.Integer))).then(function ($result3) {
                return $promise.resolve($result2 && $result3.$wrapped).then(function ($result1) {
                  return $promise.resolve($result1 && !parsed.AnotherField.$wrapped).then(function ($result0) {
                    $result = $t.fastbox($result0 && parsed.SomeInstance.AnotherBool.$wrapped, $g.____testlib.basictypes.Boolean);
                    $current = 4;
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
