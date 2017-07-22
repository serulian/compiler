$module('attributes', function () {
  var $static = this;
  $static.SimpleFunction = function (props) {
    return $t.fastbox($g.________testlib.basictypes.String.$equals($t.cast(props.$index($t.fastbox("data-foo", $g.________testlib.basictypes.String)), $g.________testlib.basictypes.String, false), $t.fastbox("bar", $g.________testlib.basictypes.String)).$wrapped && $t.cast(props.$index($t.fastbox("an-attr-here", $g.________testlib.basictypes.String)), $g.________testlib.basictypes.Boolean, false).$wrapped, $g.________testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.attributes.SimpleFunction($g.________testlib.basictypes.Mapping($t.any).overObject({
      "an-attr-here": $t.fastbox(true, $g.________testlib.basictypes.Boolean),
      "data-foo": $t.fastbox("bar", $g.________testlib.basictypes.String),
    }));
  };
});
