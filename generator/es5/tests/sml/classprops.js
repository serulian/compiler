$module('classprops', function () {
  var $static = this;
  this.$class('c89a612d', 'SomeProps', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (BoolValue, StringValue) {
      var instance = new $static();
      instance.BoolValue = BoolValue;
      instance.StringValue = StringValue;
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.SimpleFunction = function (props) {
    return $t.fastbox(($g.____testlib.basictypes.String.$equals(props.StringValue, $t.fastbox("hello world", $g.____testlib.basictypes.String)).$wrapped && props.BoolValue.$wrapped) && !(props.OptionalValue == null), $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    var $temp0;
    return $g.classprops.SimpleFunction(($temp0 = $g.classprops.SomeProps.new($t.fastbox(true, $g.____testlib.basictypes.Boolean), $t.fastbox("hello world", $g.____testlib.basictypes.String)), $temp0.OptionalValue = $t.fastbox(42, $g.____testlib.basictypes.Integer), $temp0));
  };
});
