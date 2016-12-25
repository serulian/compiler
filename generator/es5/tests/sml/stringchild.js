$module('stringchild', function () {
  var $static = this;
  $static.SimpleFunction = function (props, child) {
    return $t.fastbox(child.Length().$wrapped == 5, $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.stringchild.SimpleFunction($g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).Empty(), $t.fastbox("hello", $g.____testlib.basictypes.String));
  };
});
