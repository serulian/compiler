$module('native', function () {
  var $static = this;
  $static.MakeSomePromise = function () {
    return $t.nativenew($global.Promise)(function (r1, r2) {
      r1.call(null, $t.fastbox(42, $g.________testlib.basictypes.Integer));
      return;
    });
  };
  $static.DoSomething = $t.markpromising(function (p) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $promise.translate(p).then(function ($result0) {
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
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $promise.maybe($g.native.DoSomething($t.fastbox($g.native.MakeSomePromise(), $g.________testlib.basictypes.Promise($g.________testlib.basictypes.Integer)))).then(function ($result0) {
              $result = $t.fastbox($result0.$wrapped == 42, $g.________testlib.basictypes.Boolean);
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
  });
});
