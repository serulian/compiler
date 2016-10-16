$module('with', function () {
  var $static = this;
  this.$class('SomeReleasable', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.Release = function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $g.with.someBool = $t.box(true, $g.____testlib.basictypes.Boolean);
        $resolve();
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      return $t.createtypesig(['Release', 2, $g.____testlib.basictypes.Function($t.void).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.with.SomeReleasable).$typeref()]);
    };
  });

  $static.TEST = function () {
    var $result;
    var $temp0;
    var $current = 0;
    var $resources = $t.resourcehandler();
    var $continue = function ($resolve, $reject) {
      $resolve = $resources.bind($resolve);
      $reject = $resources.bind($reject);
      while (true) {
        switch ($current) {
          case 0:
            $t.box(123, $g.____testlib.basictypes.Integer);
            $g.with.SomeReleasable.new().then(function ($result0) {
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
            $temp0 = $result;
            $resources.pushr($temp0, '$temp0');
            $t.box(456, $g.____testlib.basictypes.Integer);
            $resources.popr('$temp0').then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $t.box(789, $g.____testlib.basictypes.Integer);
            $resolve($g.with.someBool);
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
    return $promise.new(function (resolve) {
      $static.someBool = $t.box(false, $g.____testlib.basictypes.Boolean);
      resolve();
    });
  }, '19a53bdb', []);
});
