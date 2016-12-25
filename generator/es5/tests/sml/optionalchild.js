$module('optionalchild', function () {
  var $static = this;
  $static.SimpleFunction = function (props, child) {
    return $t.fastbox(child == null, $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.optionalchild.SimpleFunction($g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).Empty());
  };
});
