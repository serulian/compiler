$module('structprops', function () {
  var $static = this;
  this.$struct('3f6b14fb', 'SomeProps', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (BoolValue, StringValue) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        BoolValue: BoolValue,
        StringValue: StringValue,
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
    $t.defineStructField($static, 'StringValue', 'StringValue', function () {
      return $g.________testlib.basictypes.String;
    }, function () {
      return $g.________testlib.basictypes.String;
    }, false);
    $t.defineStructField($static, 'OptionalValue', 'OptionalValue', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, true);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|0b2e6e78<3f6b14fb>": true,
        "equals|4|0b2e6e78<5e61c39d>": true,
        "Stringify|2|0b2e6e78<c509e19d>": true,
        "Mapping|2|0b2e6e78<204295f9<any>>": true,
        "Clone|2|0b2e6e78<3f6b14fb>": true,
        "String|2|0b2e6e78<c509e19d>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.SimpleFunction = function (props) {
    return $t.fastbox(($g.________testlib.basictypes.String.$equals(props.StringValue, $t.fastbox("hello world", $g.________testlib.basictypes.String)).$wrapped && props.BoolValue.$wrapped) && !(props.OptionalValue == null), $g.________testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    var $temp0;
    return $g.structprops.SimpleFunction(($temp0 = $g.structprops.SomeProps.new($t.fastbox(true, $g.________testlib.basictypes.Boolean), $t.fastbox("hello world", $g.________testlib.basictypes.String)), $temp0.OptionalValue = $t.fastbox(42, $g.________testlib.basictypes.Integer), $temp0));
  };
});
