$module('withas', function () {
  var $static = this;
  $static.DoSomething = function (someExpr) {
    var someName;
    var $resources = $t.resourcehandler();
    $t.fastbox(123, $g.________testlib.basictypes.Integer);
    someName = someExpr;
    $resources.pushr(someName, 'someName');
    $t.fastbox(456, $g.________testlib.basictypes.Integer);
    $resources.popr('someName');
    $t.fastbox(789, $g.________testlib.basictypes.Integer);
    var $pat = undefined;
    $resources.popall();
    return $pat;
  };
});
