$module('simple', function () {
  var $static = this;
  $static.TEST = function () {
    var a;
    a = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    return a;
  };
});
