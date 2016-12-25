$module('mappingprops', function () {
  var $static = this;
  $static.SimpleFunction = function (props) {
    return $g.____testlib.basictypes.String.$equals($t.syncnullcompare(props.$index($t.fastbox('SomeProp', $g.____testlib.basictypes.String)), function () {
      return $t.fastbox('', $g.____testlib.basictypes.String);
    }), $t.fastbox('hello world', $g.____testlib.basictypes.String));
  };
  $static.TEST = function () {
    return $g.mappingprops.SimpleFunction($g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).overObject({
      SomeProp: $t.fastbox("hello world", $g.____testlib.basictypes.String),
    }));
  };
});
