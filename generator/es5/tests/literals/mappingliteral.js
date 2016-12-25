$module('mappingliteral', function () {
  var $static = this;
  $static.TEST = function () {
    return $g.____testlib.basictypes.Mapping($g.____testlib.basictypes.Boolean).overObject(function () {
      var obj = {
      };
      obj['somekey'] = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
      obj['anotherkey'] = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
      return obj;
    }()).$index($t.fastbox('somekey', $g.____testlib.basictypes.String));
  };
});
