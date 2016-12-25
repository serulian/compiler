$module('castfunction', function () {
  var $static = this;
  $static.TEST = function () {
    $t.cast($g.castfunction.TEST, $g.____testlib.basictypes.Function($t.any), false);
    return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
  };
});
