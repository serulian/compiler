$module('mappingprops', function () {
  var $static = this;
  $static.SimpleFunction = function (props) {
    return $g.________testlib.basictypes.String.$equals($t.syncnullcompare(props.$index($t.fastbox('SomeProp', $g.________testlib.basictypes.String)), function () {
      return $t.fastbox('', $g.________testlib.basictypes.String);
    }), $t.fastbox('hello world', $g.________testlib.basictypes.String));
  };
  $static.TEST = function () {
    return $g.mappingprops.SimpleFunction($g.________testlib.basictypes.Mapping($g.________testlib.basictypes.String).overObject({
      SomeProp: $t.fastbox("hello world", $g.________testlib.basictypes.String),
    }));
  };
});
