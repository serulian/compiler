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
        "Parse|1|cf412abd<3f6b14fb>": true,
        "equals|4|cf412abd<aa28dc2d>": true,
        "Stringify|2|cf412abd<cb470bcc>": true,
        "Mapping|2|cf412abd<899aec48<any>>": true,
        "Clone|2|cf412abd<3f6b14fb>": true,
        "String|2|cf412abd<cb470bcc>": true,
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
