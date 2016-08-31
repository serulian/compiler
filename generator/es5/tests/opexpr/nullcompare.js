$module('nullcompare', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var someBool;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            someBool = null;
            $promise.resolve(someBool).then(function ($result0) {
              $result = $t.nullcompare($result0, $t.box(true, $g.____testlib.basictypes.Boolean));
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
});
