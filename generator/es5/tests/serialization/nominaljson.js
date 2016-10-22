$module('nominaljson', function () {
  var $static = this;
  this.$type('4a09c49b', 'SomeNominal', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.nominaljson.AnotherStruct;
    };
    $instance.GetValue = function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box($this, $g.nominaljson.AnotherStruct).AnotherBool);
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "GetValue|2|29dc432d<5ab5941e>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('c267cdf9', 'AnotherStruct', false, '', function () {
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
        "Parse|1|29dc432d<c267cdf9>": true,
        "equals|4|29dc432d<5ab5941e>": true,
        "Stringify|2|29dc432d<538656f2>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<c267cdf9>": true,
        "String|2|29dc432d<538656f2>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('2904f0f0', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Nested) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        Nested: Nested,
      };
      instance.$markruntimecreated();
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'Nested', 'Nested', function () {
      return $g.nominaljson.SomeNominal;
    }, function () {
      return $g.nominaljson.SomeNominal;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<2904f0f0>": true,
        "equals|4|29dc432d<5ab5941e>": true,
        "Stringify|2|29dc432d<538656f2>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<2904f0f0>": true,
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
            $g.nominaljson.AnotherStruct.new($t.fastbox(true, $g.____testlib.basictypes.Boolean)).then(function ($result1) {
              return $g.nominaljson.SomeStruct.new($t.box($result1, $g.nominaljson.SomeNominal)).then(function ($result0) {
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
            jsonString = $t.fastbox('{"Nested":{"AnotherBool":true}}', $g.____testlib.basictypes.String);
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
            $g.nominaljson.SomeStruct.Parse($g.____testlib.basictypes.JSON)(jsonString).then(function ($result0) {
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
            $promise.resolve(correct.$wrapped).then(function ($result0) {
              return ($promise.shortcircuit($result0, true) || parsed.Nested.GetValue()).then(function ($result1) {
                $result = $t.fastbox($result0 && $result1.$wrapped, $g.____testlib.basictypes.Boolean);
                $current = 4;
                $continue($resolve, $reject);
                return;
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
