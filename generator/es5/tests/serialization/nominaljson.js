$module('nominaljson', function () {
  var $static = this;
  this.$type('SomeNominal', false, '', function () {
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
      return $t.createtypesig(['GetValue', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()]);
    };
  });

  this.$struct('AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherBool) {
      var instance = new $static();
      instance.$unboxed = false;
      instance[BOXED_DATA_PROPERTY] = {
        AnotherBool: AnotherBool,
      };
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'AnotherBool', 'AnotherBool', function () {
      return $g.____testlib.basictypes.Boolean;
    }, true, function () {
      return $g.____testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.nominaljson.AnotherStruct).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.nominaljson.AnotherStruct).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.nominaljson.AnotherStruct).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
  });

  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Nested) {
      var instance = new $static();
      instance.$unboxed = false;
      instance[BOXED_DATA_PROPERTY] = {
        Nested: Nested,
      };
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'Nested', 'Nested', function () {
      return $g.nominaljson.SomeNominal;
    }, true, function () {
      return $g.nominaljson.SomeNominal;
    }, false);
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.nominaljson.SomeStruct).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.nominaljson.SomeStruct).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.nominaljson.SomeStruct).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
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
            $g.nominaljson.AnotherStruct.new($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result1) {
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
            jsonString = $t.box('{"Nested":{"AnotherBool":true}}', $g.____testlib.basictypes.String);
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
            $promise.resolve($t.unbox(correct)).then(function ($result0) {
              return ($promise.shortcircuit($result0, true) || parsed.Nested.GetValue()).then(function ($result1) {
                $result = $t.box($result0 && $t.unbox($result1), $g.____testlib.basictypes.Boolean);
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
