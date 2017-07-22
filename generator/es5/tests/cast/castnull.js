$module('castnull', function () {
  var $static = this;
  $static.TEST = function () {
    var value;
    value = null;
    return $t.cast(value, $g.________testlib.basictypes.Boolean, false);
  };
});
