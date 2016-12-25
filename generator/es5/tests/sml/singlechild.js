$module('singlechild', function () {
  var $static = this;
  $static.SimpleFunction = function (props, child) {
    return $t.fastbox(child.$wrapped == 42, $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.singlechild.SimpleFunction($g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).Empty(), $t.fastbox(42, $g.____testlib.basictypes.Integer));
  };
});
