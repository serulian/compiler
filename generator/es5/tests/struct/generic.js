$module('generic', function () {
  var $static = this;
  this.$struct('6b885b52', 'AnotherStruct', false, '', function () {
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
      return $g.________testlib.basictypes.Boolean;
    }, function () {
      return $g.________testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|0b2e6e78<6b885b52>": true,
        "equals|4|0b2e6e78<5e61c39d>": true,
        "Stringify|2|0b2e6e78<c509e19d>": true,
        "Mapping|2|0b2e6e78<204295f9<any>>": true,
        "Clone|2|0b2e6e78<6b885b52>": true,
        "String|2|0b2e6e78<c509e19d>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('d34106e1', 'SomeStruct', true, '', function (T) {
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
        "equals|4|0b2e6e78<5e61c39d>": true,
        "Stringify|2|0b2e6e78<c509e19d>": true,
        "Mapping|2|0b2e6e78<204295f9<any>>": true,
        "String|2|0b2e6e78<c509e19d>": true,
      };
      computed[("Parse|1|0b2e6e78<d34106e1<" + $t.typeid(T)) + ">>"] = true;
      computed[("Clone|2|0b2e6e78<d34106e1<" + $t.typeid(T)) + ">>"] = true;
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var ss;
    var ss2;
    ss = $g.generic.SomeStruct($g.generic.AnotherStruct).new($g.generic.AnotherStruct.new($t.fastbox(true, $g.________testlib.basictypes.Boolean)));
    ss2 = $g.generic.SomeStruct($g.________testlib.basictypes.Boolean).new($t.fastbox(true, $g.________testlib.basictypes.Boolean));
    return $t.fastbox(ss.SomeField.BoolValue.$wrapped && ss2.SomeField.$wrapped, $g.________testlib.basictypes.Boolean);
  };
});
