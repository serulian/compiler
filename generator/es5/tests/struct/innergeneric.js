$module('innergeneric', function () {
  var $static = this;
  this.$struct('078feecd', 'AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (BoolValue) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        BoolValue: BoolValue,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'BoolValue', 'BoolValue', function () {
      return $g.____testlib.basictypes.Boolean;
    }, function () {
      return $g.____testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|89b8f38e<078feecd>": true,
        "equals|4|89b8f38e<f7f23c49>": true,
        "Stringify|2|89b8f38e<549fbddd>": true,
        "Mapping|2|89b8f38e<ad6de9ce<any>>": true,
        "Clone|2|89b8f38e<078feecd>": true,
        "String|2|89b8f38e<549fbddd>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('8369fb45', 'SomeStruct', true, '', function (T) {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        SomeField: SomeField,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'SomeField', function () {
      return T;
    }, function () {
      return T;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "equals|4|89b8f38e<f7f23c49>": true,
        "Stringify|2|89b8f38e<549fbddd>": true,
        "Mapping|2|89b8f38e<ad6de9ce<any>>": true,
        "String|2|89b8f38e<549fbddd>": true,
      };
      computed[("Parse|1|89b8f38e<8369fb45<" + $t.typeid(T)) + ">>"] = true;
      computed[("Clone|2|89b8f38e<8369fb45<" + $t.typeid(T)) + ">>"] = true;
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = $t.markpromising(function () {
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
            ss = $g.innergeneric.SomeStruct($t.struct).new($g.innergeneric.AnotherStruct.new($t.fastbox(true, $g.____testlib.basictypes.Boolean)));
            $promise.maybe(ss.Stringify($g.____testlib.basictypes.JSON)()).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            jsonString = $result;
            $promise.maybe($g.innergeneric.SomeStruct($g.innergeneric.AnotherStruct).Parse($g.____testlib.basictypes.JSON)(jsonString)).then(function ($result0) {
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
            sscopy = $result;
            $promise.maybe($g.innergeneric.SomeStruct($t.struct).Parse($g.____testlib.basictypes.JSON)(jsonString)).then(function ($result0) {
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
            iss = $result;
            $resolve($t.fastbox(($t.cast(ss.SomeField, $g.innergeneric.AnotherStruct, false).BoolValue.$wrapped && sscopy.SomeField.BoolValue.$wrapped) && $t.cast(iss.SomeField, $g.innergeneric.AnotherStruct, false).BoolValue.$wrapped, $g.____testlib.basictypes.Boolean));
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
