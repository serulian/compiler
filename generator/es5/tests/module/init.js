$module('init', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.resolve($g.init.sc.value.$wrapped == 2).then(function ($result0) {
              $result = $t.fastbox($result0 && ($g.init.sc2.value.$wrapped == 4), $g.____testlib.basictypes.Boolean);
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  this.$init(function () {
    return $g.other.SomeClass.NewThing($t.fastbox(1, $g.____testlib.basictypes.Integer)).then(function ($result0) {
      $static.sc = $result0;
    });
  }, '194711ba', ['cc19450d']);
  this.$init(function () {
    return $g.other.SomeClass.NewThing($t.fastbox(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
      $static.sc2 = $result0;
    });
  }, '2e80db22', ['cc19450d']);
});
