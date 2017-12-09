$module('clone', function () {
  var $static = this;
  this.$struct('be48fdbe', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField, AnotherField) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        SomeField: SomeField,
        AnotherField: AnotherField,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'SomeField', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'AnotherField', 'AnotherField', function () {
      return $g.________testlib.basictypes.Boolean;
    }, function () {
      return $g.________testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|0b2e6e78<be48fdbe>": true,
        "equals|4|0b2e6e78<5e61c39d>": true,
        "Stringify|2|0b2e6e78<c509e19d>": true,
        "Mapping|2|0b2e6e78<204295f9<any>>": true,
        "Clone|2|0b2e6e78<be48fdbe>": true,
        "String|2|0b2e6e78<c509e19d>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $temp0;
    var first;
    var second;
    first = $g.clone.SomeStruct.new($t.fastbox(42, $g.________testlib.basictypes.Integer), $t.fastbox(false, $g.________testlib.basictypes.Boolean));
    second = ($temp0 = first.Clone(), $temp0.AnotherField = $t.fastbox(true, $g.________testlib.basictypes.Boolean), $temp0);
    return $t.fastbox(((second.SomeField.$wrapped == 42) && second.AnotherField.$wrapped) && !first.AnotherField.$wrapped, $g.________testlib.basictypes.Boolean);
  };
});
