$module('optionalchild', function () {
  var $static = this;
  $static.SimpleFunction = function (props, child) {
    return $t.fastbox(child == null, $g.________testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.optionalchild.SimpleFunction($g.________testlib.basictypes.Mapping($g.________testlib.basictypes.String).Empty());
  };
});
