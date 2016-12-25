$module('attributes', function () {
  var $static = this;
  $static.SimpleFunction = function (props) {
    return $t.fastbox($g.____testlib.basictypes.String.$equals($t.cast(props.$index($t.fastbox("data-foo", $g.____testlib.basictypes.String)), $g.____testlib.basictypes.String, false), $t.fastbox("bar", $g.____testlib.basictypes.String)).$wrapped && $t.cast(props.$index($t.fastbox("an-attr-here", $g.____testlib.basictypes.String)), $g.____testlib.basictypes.Boolean, false).$wrapped, $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.attributes.SimpleFunction($g.____testlib.basictypes.Mapping($t.any).overObject({
      "an-attr-here": $t.fastbox(true, $g.____testlib.basictypes.Boolean),
      "data-foo": $t.fastbox("bar", $g.____testlib.basictypes.String),
    }));
  };
});
