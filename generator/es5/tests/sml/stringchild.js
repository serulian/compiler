$module('stringchild', function () {
  var $static = this;
  $static.SimpleFunction = function (props, child) {
    return $t.fastbox(child.Length().$wrapped == 5, $g.________testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.stringchild.SimpleFunction($g.________testlib.basictypes.Mapping($g.________testlib.basictypes.String).Empty(), $t.fastbox("hello", $g.________testlib.basictypes.String));
  };
});
