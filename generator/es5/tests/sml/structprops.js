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
      return $g.____testlib.basictypes.Boolean;
    }, function () {
      return $g.____testlib.basictypes.Boolean;
    }, false);
    $t.defineStructField($static, 'StringValue', 'StringValue', function () {
      return $g.____testlib.basictypes.String;
    }, function () {
      return $g.____testlib.basictypes.String;
    }, false);
    $t.defineStructField($static, 'OptionalValue', 'OptionalValue', function () {
      return $g.____testlib.basictypes.Integer;
    }, function () {
      return $g.____testlib.basictypes.Integer;
    }, true);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|89b8f38e<3f6b14fb>": true,
        "equals|4|89b8f38e<f7f23c49>": true,
        "Stringify|2|89b8f38e<549fbddd>": true,
        "Mapping|2|89b8f38e<ad6de9ce<any>>": true,
        "Clone|2|89b8f38e<3f6b14fb>": true,
        "String|2|89b8f38e<549fbddd>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.SimpleFunction = function (props) {
    return $t.fastbox(($g.____testlib.basictypes.String.$equals(props.StringValue, $t.fastbox("hello world", $g.____testlib.basictypes.String)).$wrapped && props.BoolValue.$wrapped) && !(props.OptionalValue == null), $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    var $temp0;
    return $g.structprops.SimpleFunction(($temp0 = $g.structprops.SomeProps.new($t.fastbox(true, $g.____testlib.basictypes.Boolean), $t.fastbox("hello world", $g.____testlib.basictypes.String)), $temp0.OptionalValue = $t.fastbox(42, $g.____testlib.basictypes.Integer), $temp0));
  };
});
