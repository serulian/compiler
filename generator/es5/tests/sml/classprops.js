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
    return $t.fastbox(($g.________testlib.basictypes.String.$equals(props.StringValue, $t.fastbox("hello world", $g.________testlib.basictypes.String)).$wrapped && props.BoolValue.$wrapped) && !(props.OptionalValue == null), $g.________testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    var $temp0;
    return $g.classprops.SimpleFunction(($temp0 = $g.classprops.SomeProps.new($t.fastbox(true, $g.________testlib.basictypes.Boolean), $t.fastbox("hello world", $g.________testlib.basictypes.String)), $temp0.OptionalValue = $t.fastbox(42, $g.________testlib.basictypes.Integer), $temp0));
  };
});
