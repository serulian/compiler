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
        "Parse|1|29dc432d<3f6b14fb>": true,
        "equals|4|29dc432d<43834c3f>": true,
        "Stringify|2|29dc432d<5cffd9b5>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<3f6b14fb>": true,
        "String|2|29dc432d<5cffd9b5>": true,
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
