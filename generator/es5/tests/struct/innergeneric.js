$module('innergeneric', function () {
  var $static = this;
  this.$struct('AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (BoolValue) {
      var instance = new $static();
      var init = [];
      instance.$unboxed = false;
      instance[BOXED_DATA_PROPERTY] = {
        BoolValue: BoolValue,
      };
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.$fields = [];
    $t.defineStructField($static, 'BoolValue', 'BoolValue', function () {
      return $g.____testlib.basictypes.Boolean;
    }, true, function () {
      return $g.____testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.innergeneric.AnotherStruct).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.innergeneric.AnotherStruct).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.innergeneric.AnotherStruct).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
  });

  this.$struct('SomeStruct', true, '', function (T) {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      var init = [];
      instance.$unboxed = false;
      instance[BOXED_DATA_PROPERTY] = {
        SomeField: SomeField,
      };
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'SomeField', function () {
      return T;
    }, false, function () {
      return T;
    }, false);
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.innergeneric.SomeStruct(T)).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.innergeneric.SomeStruct(T)).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.innergeneric.SomeStruct(T)).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
  });

  $static.TEST = function () {
    var $result;
    var iss;
    var jsonString;
    var ss;
    var sscopy;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.innergeneric.AnotherStruct.new($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result1) {
              return $g.innergeneric.SomeStruct($t.struct).new($result1).then(function ($result0) {
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
            ss = $result;
            ss.Stringify($g.____testlib.basictypes.JSON)().then(function ($result0) {
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
            jsonString = $result;
            $g.innergeneric.SomeStruct($g.innergeneric.AnotherStruct).Parse($g.____testlib.basictypes.JSON)(jsonString).then(function ($result0) {
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
            sscopy = $result;
            $g.innergeneric.SomeStruct($t.struct).Parse($g.____testlib.basictypes.JSON)(jsonString).then(function ($result0) {
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
            iss = $result;
            $promise.resolve($t.unbox($t.cast(ss.SomeField, $g.innergeneric.AnotherStruct, false).BoolValue)).then(function ($result1) {
              return $promise.resolve($result1 && $t.unbox(sscopy.SomeField.BoolValue)).then(function ($result0) {
                $result = $t.box($result0 && $t.unbox($t.cast(iss.SomeField, $g.innergeneric.AnotherStruct, false).BoolValue), $g.____testlib.basictypes.Boolean);
                $current = 5;
                $continue($resolve, $reject);
                return;
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
