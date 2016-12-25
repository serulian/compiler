$module('assertnotnull', function () {
  var $static = this;
  $static.TEST = function () {
    var someValue;
    someValue = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    return $t.assertnotnull(someValue);
  };
});
