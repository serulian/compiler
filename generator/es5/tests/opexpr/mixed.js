$module('mixed', function () {
  var $static = this;
  $static.TEST = function () {
    var finalIndex;
    finalIndex = $t.fastbox(-2, $g.____testlib.basictypes.Integer);
    return $t.fastbox((finalIndex.$wrapped >= 10) || (finalIndex.$wrapped < 0), $g.____testlib.basictypes.Boolean);
  };
});
