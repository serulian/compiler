$module('singlechild', function () {
  var $static = this;
  $static.SimpleFunction = function (props, child) {
    return $t.fastbox(child.$wrapped == 42, $g.________testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.singlechild.SimpleFunction($g.________testlib.basictypes.Mapping($g.________testlib.basictypes.String).Empty(), $t.fastbox(42, $g.________testlib.basictypes.Integer));
  };
});
