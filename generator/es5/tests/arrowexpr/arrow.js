$module('arrow', function () {
  var $static = this;
  this.$class('SomePromise', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Then = function (resolve) {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              resolve($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
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
              $result;
              $resolve($this);
              return;

            default:
              $resolve();
              return;
          }
        }
      };
      return $promise.new($continue);
    };
    $instance.Catch = function (rejection) {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($this);
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      return $t.createtypesig(['Then', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Awaitable($g.____testlib.basictypes.Boolean)).$typeref()], ['Catch', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Awaitable($g.____testlib.basictypes.Boolean)).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.arrow.SomePromise).$typeref()]);
    };
  });

  $static.DoSomething = function (p) {
    var somebool;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.translate(p).then(function (resolved) {
              somebool = resolved;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (rejected) {
              $reject(rejected);
              return;
            });
            return;

          case 1:
            $resolve(somebool);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.arrow.SomePromise.new().then(function ($result1) {
              return $g.arrow.DoSomething($result1).then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
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
