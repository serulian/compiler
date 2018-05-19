$module('inline', function () {
  var $static = this;
  $static.TEST = $t.markpromising(function () {
    var $result;
    var result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            result = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
            $promise.maybe((function () {
              result = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
              return;
            })()).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve(result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
